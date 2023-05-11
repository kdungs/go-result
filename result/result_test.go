package result_test

import (
	"errors"
	"testing"

	"github.com/kdungs/go-result/result"
)

func FuzzWrapUnwrap(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string, isErr bool) {
		var v string
		var err error
		if isErr {
			v = s
			err = nil
		} else {
			v = *new(string)
			err = errors.New(s)
		}
		r := result.Wrap(v, err)
		uv, uerr := r.Unwrap()
		if err != uerr {
			t.Fatalf("want %v, got %v", err, uerr)
		}
		if err == nil && v != uv {
			t.Fatalf("want %q, got %q", v, uv)
		}
	})
}

func FuzzWrapUnwrapPointer(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string, isErr bool) {
		var v string
		var err error
		if isErr {
			v = s
			err = nil
		} else {
			v = *new(string)
			err = errors.New(s)
		}
		r := result.Wrap(&v, err)
		uv, uerr := r.Unwrap()
		if err != uerr {
			t.Fatalf("want %v, got %v", err, uerr)
		}
		if err == nil && v != *uv {
			t.Fatalf("want %q, got %q", v, *uv)
		}
	})
}

func FuzzOf(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		r := result.Of(s)
		v, err := r.Unwrap()
		if err != nil {
			t.Fatalf("want no error, got %v", err)
		}
		if v != s {
			t.Fatalf("want %q, got %q", s, v)
		}
	})
}

func FuzzOfPointer(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		r := result.Of(&s)
		v, err := r.Unwrap()
		if err != nil {
			t.Fatalf("want no error, got %v", err)
		}
		if *v != s {
			t.Fatalf("want %q, got %q", s, *v)
		}
	})
}

func TestOr(t *testing.T) {
	cases := []struct {
		name     string
		r        result.R[string]
		d        string
		expected string
	}{
		{
			name:     "no value uses default",
			r:        result.OfErr[string](errors.New("missing value")),
			d:        "default value",
			expected: "default value",
		},
		{
			name:     "uses value if present",
			r:        result.Of("foo bar"),
			d:        "bar baz",
			expected: "foo bar",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v := tc.r.Or(tc.d)
			if v != tc.expected {
				t.Fatalf("want %q, got %q", tc.expected, v)
			}
		})
	}
}

func TestOfErrPointer(t *testing.T) {
	errExample := errors.New("test")
	r := result.OfErr[*int](errExample)
	if _, err := r.Unwrap(); err != errExample {
		t.Fatalf("want %v, got %v", errExample, err)
	}
}
