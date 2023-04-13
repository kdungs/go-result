package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdungs/go-result"
)

func TestFmap(t *testing.T) {
	errV := errors.New("v")
	suites := []struct {
		name string
		f    func(func(int) string, result.R[int]) result.R[string]
	}{
		{
			name: "Fmap",
			f: func(f func(int) string, r result.R[int]) result.R[string] {
				return result.Fmap(f)(r)
			},
		},
		{
			name: "EagerFmap",
			f:    result.EagerFmap[int, string],
		},
	}
	cases := []struct {
		name        string
		v           result.R[int]
		expectedErr error
		expectedVal string
	}{
		{
			name:        "error",
			v:           result.OfErr[int](errV),
			expectedErr: errV,
		},
		{
			name:        "value",
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
					r := ts.f(func(x int) string { return fmt.Sprint(x) }, tc.v)
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
