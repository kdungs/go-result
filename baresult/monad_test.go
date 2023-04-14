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
	providingError := func() (int, error) {
		return 0, errV
	}
	providingValue := func() (int, error) {
		return 42, nil
	}

	cases := []struct {
		name        string
		f           func(int) (string, error)
		provideV    func() (int, error)
		expectedErr error
		expectedVal string
	}{
		{
			name:        "both are error",
			f:           returningError,
			provideV:    providingError,
			expectedErr: errV,
		},
		{
			name:        "f is error",
			f:           returningError,
			provideV:    providingValue,
			expectedErr: errF,
		},
		{
			name:        "v is error",
			f:           returningValue,
			provideV:    providingError,
			expectedErr: errV,
		},
		{
			name:        "both are value",
			f:           returningValue,
			provideV:    providingValue,
			expectedErr: nil,
			expectedVal: "42",
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v, err := baresult.Bind(tc.f)(tc.provideV())
			if err != tc.expectedErr {
				t.Fatalf("want %v, got %v", tc.expectedErr, err)
			}
			if tc.expectedErr == nil && v != tc.expectedVal {
				t.Fatalf("want %q, got %q", tc.expectedVal, v)
			}
		})
	}
}
