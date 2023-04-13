package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdungs/go-result"
)

func TestKleisli(t *testing.T) {
	errF := errors.New("a")
	fErr := func(int) result.R[string] {
		return result.OfErr[string](errF)
	}
	fVal := func(x int) result.R[string] {
		return result.Of(fmt.Sprint(x))
	}
	errG := errors.New("g")
	gErr := func(string) result.R[string] {
		return result.OfErr[string](errG)
	}
	gVal := func(s string) result.R[string] {
		return result.Of(fmt.Sprintf("[%s]", s))
	}
	suites := []struct {
		name string
		f    func(func(int) result.R[string], func(string) result.R[string], int) result.R[string]
	}{
		{
			name: "Kleisli",
			f: func(f func(int) result.R[string], g func(string) result.R[string], v int) result.R[string] {
				return result.Kleisli(f, g)(v)
			},
		},
		{
			name: "EagerKleisli",
			f:    result.EagerKleisli[int, string, string],
		},
	}
	cases := []struct {
		name        string
		f           func(int) result.R[string]
		g           func(string) result.R[string]
		v           int
		expectedErr error
		expectedVal string
	}{
		{
			name:        "both return error",
			f:           fErr,
			g:           gErr,
			expectedErr: errF,
		},
		{
			name:        "f returns error",
			f:           fErr,
			g:           gVal,
			expectedErr: errF,
		},
		{
			name:        "g returns error",
			f:           fVal,
			g:           gErr,
			expectedErr: errG,
		},
		{
			name:        "both return value",
			f:           fVal,
			g:           gVal,
			v:           42,
			expectedErr: nil,
			expectedVal: "[42]",
		},
	}
	for _, ts := range suites {
		ts := ts
		t.Run(ts.name, func(t *testing.T) {
			for _, tc := range cases {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					r := ts.f(tc.f, tc.g, tc.v)
					v, err := r.Unwrap()
					if err != tc.expectedErr {
						t.Fatalf("want %v, got %v", tc.expectedErr, err)
					}
					if tc.expectedErr == nil && v != tc.expectedVal {
						t.Fatalf("want %q, got %q", tc.expectedVal, v)
					}
				})
			}
		})
	}
}
