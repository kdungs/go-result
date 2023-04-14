package baresult_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdungs/go-result/baresult"
)

func TestZipWith(t *testing.T) {
	errA := errors.New("a")
	errB := errors.New("b")
	cases := []struct {
		name        string
		provideA    func() (int, error)
		provideB    func() (string, error)
		expectedErr error
		expectedVal string
	}{
		{
			name:        "both are error",
			provideA:    func() (int, error) { return 0, errA },
			provideB:    func() (string, error) { return "", errB },
			expectedErr: errA,
		},
		{
			name:        "a is error",
			provideA:    func() (int, error) { return 0, errA },
			provideB:    func() (string, error) { return "foo", nil },
			expectedErr: errA,
		},
		{
			name:        "b is error",
			provideA:    func() (int, error) { return 42, nil },
			provideB:    func() (string, error) { return "", errB },
			expectedErr: errB,
		},
		{
			name:        "both are value",
			provideA:    func() (int, error) { return 42, nil },
			provideB:    func() (string, error) { return "foo", nil },
			expectedErr: nil,
			expectedVal: "foo42",
		},
	}
	f := func(a int, b string) string {
		return fmt.Sprintf("%s%d", b, a)
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			v, err := baresult.ZipWith(f)(tc.provideA())(tc.provideB())
			if err != tc.expectedErr {
				t.Fatalf("want %v, got %v", tc.expectedErr, err)
			}
			if tc.expectedErr == nil && v != tc.expectedVal {
				t.Fatalf("want %q, got %q", tc.expectedVal, v)
			}
		})
	}
}
