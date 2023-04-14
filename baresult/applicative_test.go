package baresult_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdungs/go-result/baresult"
)

func TestApply(t *testing.T) {
	errF := errors.New("f")
	errV := errors.New("v")

	suites := []struct {
		name string
		f    func(func(int) string, error, int, error) (string, error)
	}{
		{
			name: "Apply",
			f: func(f func(int) string, errF error, v int, errV error) (string, error) {
				return baresult.Apply(f, errF)(v, errV)
			},
		},
		{
			name: "EagerApply",
			f:    baresult.EagerApply[int, string],
		},
	}
	cases := []struct {
		name        string
		f           func(int) string
		errF        error
		v           int
		errV        error
		expectedErr error
		expectedVal string
	}{
		{
			name:        "both are error",
			errF:        errF,
			errV:        errV,
			expectedErr: errF,
		},
		{
			name:        "f is error",
			errF:        errF,
			v:           42,
			expectedErr: errF,
		},
		{
			name:        "v is error",
			f:           func(x int) string { return "42" },
			errV:        errV,
			expectedErr: errV,
		},
		{
			name:        "both are value",
			f:           func(x int) string { return fmt.Sprint(x) },
			v:           42,
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
					v, err := ts.f(tc.f, tc.errF, tc.v, tc.errV)
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
