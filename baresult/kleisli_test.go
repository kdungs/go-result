package baresult_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdungs/go-result/baresult"
)

func TestKleisli(t *testing.T) {
	errF := errors.New("a")
	fErr := func(int) (string, error) {
		return "", errF
	}
	fVal := func(x int) (string, error) {
		return fmt.Sprint(x), nil
	}
	errG := errors.New("g")
	gErr := func(string) (string, error) {
		return "", errG
	}
	gVal := func(s string) (string, error) {
		return fmt.Sprintf("[%s]", s), nil
	}
	suites := []struct {
		name string
		f    func(func(int) (string, error), func(string) (string, error), int) (string, error)
	}{
		{
			name: "Kleisli",
			f: func(f func(int) (string, error), g func(string) (string, error), v int) (string, error) {
				return baresult.Kleisli(f, g)(v)
			},
		},
		{
			name: "EagerKleisli",
			f:    baresult.EagerKleisli[int, string, string],
		},
	}
	cases := []struct {
		name        string
		f           func(int) (string, error)
		g           func(string) (string, error)
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
					v, err := ts.f(tc.f, tc.g, tc.v)
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
