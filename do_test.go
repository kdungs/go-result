package result_test

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/kdungs/go-result"
)

var errV = errors.New("no value")

func makeV(i *int) result.R[*int]   { return result.Of(i) }
func makeErr(_ *int) result.R[*int] { return result.OfErr[*int](errV) }

func TestDoValue(t *testing.T) {
	cases := []struct {
		name        string
		i           int
		makeR       func(i *int) result.R[*int]
		expectedErr error
		expected    int
	}{
		{
			name:        "value",
			i:           3,
			makeR:       makeV,
			expectedErr: nil,
			expected:    4,
		},
		{
			name:        "error",
			makeR:       makeErr,
			expectedErr: errV,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			err := result.Do(tc.makeR(&(tc.i)), func(i *int) { *i++ })
			if err != tc.expectedErr {
				t.Fatalf("got %v, want %v", err, tc.expectedErr)
			}
			if err == nil && tc.i != tc.expected {
				t.Fatalf("got %d, want %d", tc.i, tc.expected)
			}
		})
	}
}

func TestDoE(t *testing.T) {
	errF := errors.New("f")
	type sub struct {
		name        string
		i           int
		makeR       func(*int) result.R[*int]
		expectedErr error
		expected    int
	}
	cases := []struct {
		name string
		f    func(*int) error
		subs []sub
	}{
		{
			name: "function errors",
			f:    func(_ *int) error { return errF },
			subs: []sub{
				{
					name:        "value",
					i:           23,
					makeR:       makeV,
					expectedErr: errF,
				},
				{
					name:        "error",
					makeR:       makeErr,
					expectedErr: errV,
				},
			},
		},
		{
			name: "function succeeds",
			f: func(i *int) error {
				*i++
				return nil
			},
			subs: []sub{
				{
					name:        "value",
					i:           23,
					makeR:       makeV,
					expectedErr: nil,
					expected:    24,
				},
				{
					name:        "error",
					makeR:       makeErr,
					expectedErr: errV,
				},
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for _, ts := range tc.subs {
				ts := ts
				t.Run(ts.name, func(t *testing.T) {
					err := result.DoE(ts.makeR(&ts.i), tc.f)
					if err != ts.expectedErr {
						t.Fatalf("got %v, want %v", err, ts.expectedErr)
					}
					if err == nil && ts.i != ts.expected {
						t.Fatalf("got %d, want %d", ts.i, ts.expected)
					}
				})
			}
		})
	}
}

func TestDoZip(t *testing.T) {
	f := func(a *[]int, b int) {
		*a = append(*a, b)
	}
	errA := errors.New("a")
	errB := errors.New("b")
	cases := []struct {
		name        string
		errA        error
		b           result.R[int]
		expectedErr error
		expectedLen int
	}{
		{
			name:        "both error",
			errA:        errA,
			b:           result.OfErr[int](errB),
			expectedErr: errA,
		},
		{
			name:        "a error",
			errA:        errA,
			b:           result.Of(23),
			expectedErr: errA,
		},
		{
			name:        "b error",
			errA:        nil,
			b:           result.OfErr[int](errB),
			expectedErr: errB,
		},
		{
			name:        "both value",
			errA:        nil,
			b:           result.Of(23),
			expectedErr: nil,
			expectedLen: 1,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var ints []int
			a := result.Of(&ints)
			if tc.errA != nil {
				a = result.OfErr[*[]int](tc.errA)
			}
			err := result.DoZip(a, tc.b, f)
			if err != tc.expectedErr {
				t.Fatalf("got %v, want %v", err, tc.expectedErr)
			}
			if len(ints) != tc.expectedLen {
				t.Fatalf("got len=%d, want %d", len(ints), tc.expectedLen)
			}
		})
	}
}

func TestDoZipE(t *testing.T) {
	errF := errors.New("f")
	errW := errors.New("w")
	errA := errors.New("a")
	type sub struct {
		name        string
		errW        error
		a           result.R[string]
		expectedErr error
		expected    string
	}

	cases := []struct {
		name string
		f    func(w io.Writer, a string) error
		subs []sub
	}{
		{
			name: "function errors",
			f: func(_ io.Writer, _ string) error {
				return errF
			},
			subs: []sub{
				{
					name:        "both error",
					errW:        errW,
					a:           result.OfErr[string](errA),
					expectedErr: errW,
				},
				{
					name:        "w error",
					errW:        errW,
					a:           result.Of("foo bar"),
					expectedErr: errW,
				},
				{
					name:        "a error",
					errW:        nil,
					a:           result.OfErr[string](errA),
					expectedErr: errA,
				},
				{
					name:        "both value",
					errW:        nil,
					a:           result.Of("foo bar"),
					expectedErr: errF,
				},
			},
		},
		{
			name: "function succeeds",
			f: func(w io.Writer, a string) error {
				_, err := fmt.Fprintf(w, "%s", a)
				return err
			},
			subs: []sub{
				{
					name:        "both error",
					errW:        errW,
					a:           result.OfErr[string](errA),
					expectedErr: errW,
				},
				{
					name:        "w error",
					errW:        errW,
					a:           result.Of("foo bar"),
					expectedErr: errW,
				},
				{
					name:        "a error",
					errW:        nil,
					a:           result.OfErr[string](errA),
					expectedErr: errA,
				},
				{
					name:        "both value",
					errW:        nil,
					a:           result.Of("foo bar"),
					expectedErr: nil,
					expected:    "foo bar",
				},
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			for _, ts := range tc.subs {
				ts := ts
				t.Run(ts.name, func(t *testing.T) {
					var s strings.Builder
					w := result.Of(io.Writer(&s))
					if ts.errW != nil {
						w = result.OfErr[io.Writer](ts.errW)
					}
					err := result.DoZipE(w, ts.a, tc.f)
					if err != ts.expectedErr {
						t.Fatalf("got %v, want nil", err)
					}
					if s.String() != ts.expected {
						t.Fatalf("got %q, want %q", s.String(), ts.expected)
					}
				})
			}
		})
	}

}
