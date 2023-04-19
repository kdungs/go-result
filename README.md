# go-result

_Result types for Go; because `(T, error)` can be considered a ✨ monad ✨._

> Your scientists were so preoccupied with whether or not they could that they
> didn't stop to think if they should!

**Disclaimer**: This whole thing is purely for (mostly my personal) education.
Don't try this ~~at home~~ in prod.

TIL there's actually an [ongoing discussion around how to improve error
handling in Go
2](https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling.md).

Why would this be useful? Take a look at [example_test.go](example_test.go).


## Caveats

 - No shortcuts! If an error happens early on, all computations still happen in
   the form of forwarding the error. With the current implementation this also
   means constructing a few default values on the way.
 - `defer Cleanup()` becomes a bit clunky, e.g.

```go
fh := result.Wrap(os.Open(fname))
defer result.Do(fh, func(f *os.File) { f.Close() })
```
