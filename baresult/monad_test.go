package baresult_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdungs/go-result/baresult"
)

func TestBind(t *testing.T) {
	errF := errors.New("f")
	errV := errors.New("v")
	returningError := func(int) (string, error) {
		return "", errF
	}
	returningValue := func(x int) (string, error) {
		return fmt.Sprint(x), nil
	}

	suites := []struct {
		name string
		f    func(func(int) (string, error), int, error) (string, error)
	}{
		{
			name: "Bind",
			f: func(f func(int) (string, error), v int, err error) (string, error) {
				return baresult.Bind(f)(v, err)
			},
		},
		{
			name: "EagerBind",
			f:    baresult.EagerBind[int, string],
		},
	}
	cases := []struct {
		name        string
		f           func(int) (string, error)
		v           int
		err         error
		expectedErr error
		expectedVal string
	}{
		{
			name:        "both are error",
			f:           returningError,
			err:         errV,
			expectedErr: errV,
		},
		{
			name:        "f is error",
			f:           returningError,
			v:           42,
			expectedErr: errF,
		},
		{
			name:        "v is error",
			f:           returningValue,
			err:         errV,
			expectedErr: errV,
		},
		{
			name:        "both are value",
			f:           returningValue,
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
					v, err := ts.f(tc.f, tc.v, tc.err)
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
