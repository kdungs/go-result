package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdungs/go-result"
)

func TestApply(t *testing.T) {
	errF := errors.New("f")
	errV := errors.New("v")

	suites := []struct {
		name string
		f    func(result.R[func(int) string], result.R[int]) result.R[string]
	}{
		{
			name: "Apply",
			f: func(f result.R[func(int) string], r result.R[int]) result.R[string] {
				return result.Apply(f)(r)
			},
		},
		{
			name: "EagerApply",
			f:    result.EagerApply[int, string],
		},
	}
	cases := []struct {
		name        string
		f           result.R[func(int) string]
		v           result.R[int]
		expectedErr error
		expectedVal string
	}{
		{
			name:        "both are error",
			f:           result.OfErr[func(int) string](errF),
			v:           result.OfErr[int](errV),
			expectedErr: errF,
		},
		{
			name:        "f is error",
			f:           result.OfErr[func(int) string](errF),
			v:           result.Of(42),
			expectedErr: errF,
		},
		{
			name:        "v is error",
			f:           result.Of(func(x int) string { return "42" }),
			v:           result.OfErr[int](errV),
			expectedErr: errV,
		},
		{
			name:        "both are value",
			f:           result.Of(func(x int) string { return fmt.Sprint(x) }),
			v:           result.Of(42),
			expectedErr: nil,
			expectedVal: "42",
		},
	}
	for _, ts := range suites {
		ts := ts
		t.Run(ts.name, func(t *testing.T) {
			for _, tc := range cases {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					r := ts.f(tc.f, tc.v)
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
