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
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := result.Kleisli(tc.f, tc.g)(tc.v)
			v, err := r.Unwrap()
			if err != tc.expectedErr {
				t.Fatalf("want %v, got %v", tc.expectedErr, err)
			}
			if tc.expectedErr == nil && v != tc.expectedVal {
				t.Fatalf("want %q, got %q", tc.expectedVal, v)
			}
		})
	}
}
