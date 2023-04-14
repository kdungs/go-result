package baresult_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdungs/go-result/baresult"
)

func TestFmap(t *testing.T) {
	errTest := errors.New("test")
	cases := []struct {
		name        string
		provide     func() (int, error)
		expectedErr error
		expectedVal string
	}{
		{
			name:        "error",
			provide:     func() (int, error) { return 0, errTest },
			expectedErr: errTest,
		},
		{
			name:        "value",
			provide:     func() (int, error) { return 42, nil },
			expectedErr: nil,
			expectedVal: "42",
		},
	}
	f := func(x int) string { return fmt.Sprint(x) }
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			val, err := baresult.Fmap(f)(tc.provide())
			if err != tc.expectedErr {
				t.Fatalf("want %v, got %v", tc.expectedErr, err)
			}
			if tc.expectedErr == nil && val != tc.expectedVal {
				t.Fatalf("want %q, got %q", tc.expectedVal, val)
			}
		})
	}
}
