package baresult_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdungs/go-result/baresult"
)

func TestFmap(t *testing.T) {
	errV := errors.New("v")
	suites := []struct {
		name string
		f    func(func(int) string, int, error) (string, error)
	}{
		{
			name: "Fmap",
			f: func(f func(int) string, v int, err error) (string, error) {
				return baresult.Fmap(f)(v, err)
			},
		},
		{
			name: "EagerFmap",
			f:    baresult.EagerFmap[int, string],
		},
	}
	cases := []struct {
		name        string
		val         int
		err         error
		expectedErr error
		expectedVal string
	}{
		{
			name:        "error",
			val:         0,
			err:         errV,
			expectedErr: errV,
		},
		{
			name:        "value",
			val:         42,
			err:         nil,
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
					val, err := ts.f(
						func(x int) string { return fmt.Sprint(x) },
						tc.val,
						tc.err,
					)
					if err != tc.expectedErr {
						t.Fatalf("want %v, got %v", tc.expectedErr, err)
					}
					if tc.expectedErr == nil && val != tc.expectedVal {
						t.Fatalf("want %q, got %q", tc.expectedVal, val)
					}
				})
			}
		})
	}
}
