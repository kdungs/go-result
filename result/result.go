// Package result defines a result type `R[T]` that encapsulates a value and an
// error as well as functions to work with `R[T]` conveniently.
//
// Other languages would call this `Either` (Haskell), `Result` (Rust),
// `StatusOr` (C++/Abseil), or `expected` (C++23).
//
// On top of that, the package provides a number of functions that make working
// with `R[T]` a bit more ergonomic. Also, `R[T]` is a ✨ monad ✨, wow!
package result

// R is a generic result that can either hold a value of type T or an error.
type R[T any] struct {
	err error
	v   T
}

// Unwrap returns both the value and the error component of `r`. Only one of
// them will hold a meaningful value as it is custom for Go functions that
// return a value and an error. Calling code should check `err != nil`.
func (r R[T]) Unwrap() (T, error) {
	return r.v, r.err
}

// Or returns the value part of `r` or the provided default value `d` in case
// `r` is actually holding an error.
func (r R[T]) Or(d T) T {
	if r.err != nil {
		return d
	}
	return r.v
}

// Wrap takes a value and an error (e.g. from a function returning those two)
// and turns them into an `R[T]`.
func Wrap[T any](v T, err error) R[T] {
	return R[T]{err, v}
}

// Of constructs an `R[T]` from a value of type `T`.
func Of[T any](v T) R[T] {
	return Wrap(v, nil)
}

// OfErr[T] constructs an `R[T]` holding an error.
// Since `T` cannot be inferred from the error alone, calling code is required
// to provide it, e.g. `OfErr[int](errors.New("not a positive number"))`.
func OfErr[T any](err error) R[T] {
	return Wrap(*new(T), err)
}
