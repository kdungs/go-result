package result_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/kdungs/go-result"
)

func TestZipWith(t *testing.T) {
	errA := errors.New("a")
	errB := errors.New("b")
	cases := []struct {
		name        string
		a           result.R[int]
		b           result.R[string]
		expectedErr error
		expectedVal string
	}{
		{
			name:        "both are error",
			a:           result.OfErr[int](errA),
			b:           result.OfErr[string](errB),
			expectedErr: errA,
		},
		{
			name:        "a is error",
			a:           result.OfErr[int](errA),
			b:           result.Of("foo"),
			expectedErr: errA,
		},
		{
			name:        "b is error",
			a:           result.Of(42),
			b:           result.OfErr[string](errB),
			expectedErr: errB,
		},
		{
			name:        "both are value",
			a:           result.Of(42),
			b:           result.Of("foo"),
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
			r := result.ZipWith(f)(tc.a)(tc.b)
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
