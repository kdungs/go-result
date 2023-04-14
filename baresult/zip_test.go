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
	suites := []struct {
		name string
		f    func(func(int, string) string, int, error, string, error) (string, error)
	}{
		{
			name: "ZipWith",
			f: func(
				f func(int, string) string,
				a int,
				errA error,
				b string,
				errB error,
			) (string, error) {
				return baresult.ZipWith(f)(a, errA, b, errB)
			},
		},
		{
			name: "EagerZipWith",
			f:    baresult.EagerZipWith[int, string, string],
		},
	}
	cases := []struct {
		name        string
		a           int
		errA        error
		b           string
		errB        error
		expectedErr error
		expectedVal string
	}{
		{
			name:        "both are error",
			errA:        errA,
			errB:        errB,
			expectedErr: errA,
		},
		{
			name:        "a is error",
			errA:        errA,
			b:           "foo",
			expectedErr: errA,
		},
		{
			name:        "b is error",
			a:           42,
			errB:        errB,
			expectedErr: errB,
		},
		{
			name:        "both are value",
			a:           42,
			b:           "foo",
			expectedErr: nil,
			expectedVal: "foo42",
		},
	}
	for _, ts := range suites {
		ts := ts
		t.Run(ts.name, func(t *testing.T) {
			for _, tc := range cases {
				tc := tc
				t.Run(tc.name, func(t *testing.T) {
					v, err := ts.f(
						func(a int, b string) string {
							return fmt.Sprintf("%s%d", b, a)
						},
						tc.a,
						tc.errA,
						tc.b,
						tc.errB,
					)
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
