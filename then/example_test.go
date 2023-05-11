package then_test

import (
	"fmt"
	"io"
	"log"
	"sort"
	"strings"

	"github.com/kdungs/go-result/then"
)

// Look ma, almost no `if err != nil`.
func Example() {
	// Caveat: There's no "defer Close" now.
	pipeline := then.Merge(
		fakeCreate,
		then.Map(
			then.Map(
				then.Chain(fakeOpen, io.ReadAll),
				func(bs []byte) string { return string(bs) },
			),
			countWords,
		),
		writeSortedCounts,
	)

	var buf strings.Builder
	if err := pipeline(&buf, "filename"); err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("%s", buf.String())
	// Output:
	// But: 1
	// Hello: 1
	// This: 2
	// an: 1
	// coolest: 1
	// ever: 1
	// example: 1
	// honest: 1
	// is: 3
	// it: 1
	// not: 1
	// probably: 1
	// seen.: 1
	// text.: 1
	// the: 1
	// thing: 1
	// work.: 1
	// world!: 1
	// you've: 1
}

func fakeOpen(_ string) (io.Reader, error) {
	// This is cheating a bit. Normally, if we'd use `os.Open`, we'd get an
	// `*os.File` which cannot be used in `Map` with a function that expects an
	// `io.Reader`.
	// Normally this would look like
	//  ifh := result.Wrap(os.Open("LICENSE"))
	//  in := result.Map(ifh, func(f *os.File) io.Reader {
	//  	return f
	//  })
	text := `Hello world!
This is an example text.
This is probably not the coolest thing you've ever seen.
But it is honest work.
`
	return strings.NewReader(text), nil
}

func fakeCreate(buf *strings.Builder) (io.Writer, error) {
	// Same as `fakeOpen` we're cheating by returning the interface instead of
	// `*os.File`. This is more because faking an `os.File` seems very hard...
	// Normally, this would look like
	//  ofh := result.Wrap(os.CreateTemp("", "counts"))
	//  of := result.Map(ofh, func(f *os.File) io.Writer {
	//  	return f
	//  })
	return buf, nil
}

func countWords(s string) map[string]int {
	cnts := make(map[string]int)
	for _, w := range strings.Fields(s) {
		cnts[w]++
	}
	return cnts
}

func writeSortedCounts(w io.Writer, cnts map[string]int) error {
	sortedWords := make([]string, 0, len(cnts))
	for word, _ := range cnts {
		sortedWords = append(sortedWords, word)
	}
	sort.Strings(sortedWords)
	for _, word := range sortedWords {
		if _, err := fmt.Fprintf(w, "%s: %d\n", word, cnts[word]); err != nil {
			return err
		}
	}
	return nil
}
