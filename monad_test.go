package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdungs/go-result"
)

func TestBind(t *testing.T) {
	errF := errors.New("f")
	errV := errors.New("v")
	returningError := func(int) result.R[string] {
		return result.OfErr[string](errF)
	}
	returningValue := func(x int) result.R[string] {
		return result.Of(fmt.Sprint(x))
	}
	cases := []struct {
		name        string
		f           func(int) result.R[string]
		v           result.R[int]
		expectedErr error
		expectedVal string
	}{
		{
			name:        "both are error",
			f:           returningError,
			v:           result.OfErr[int](errV),
			expectedErr: errV,
		},
		{
			name:        "f is error",
			f:           returningError,
			v:           result.Of(42),
			expectedErr: errF,
		},
		{
			name:        "v is error",
			f:           returningValue,
			v:           result.OfErr[int](errV),
			expectedErr: errV,
		},
		{
			name:        "both are value",
			f:           returningValue,
			v:           result.Of(42),
			expectedErr: nil,
			expectedVal: "42",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := result.Bind(tc.f)(tc.v)
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
