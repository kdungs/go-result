// Package result defines a result type `R[T]` that encapsulates a value and an
// error.
// Other languages would call this `Either` (Haskell), `Result` (Rust),
// `StatusOr` (C++/Abseil), or `expected` (C++23).
// On top of that, the package provides a number of functions that make working
// with `R[T]` a bit more ergonomic. Since `R[T]` is a ✨ monad ✨ (and thus also
// an applicative functor), functions like `Fmap` and `Bind` emerge trivially
// and can be used to compose chains of function calls without repeated
// `if err != nil ...`.
//
// Other than Haskell, this package defines `bind` (`>>=`) with the arguments
// flipped to be consistent with the definitions of `fmap` and `apply` (`<*>`).
// This is especially useful because Go doesn't have Haskell's facilities for
// currying and function composition and it seems more ergonomic to be able to
// "bind a function" (lift?) and then apply that to an `R[T]` rather than to
// take an `R[T]` and return a function that accepts a function that then acts
// on the value or error...
//
//	fmap  ::   (a -> b) -> R a -> R b
//	apply :: R (a -> b) -> R a -> R b
//	bind  :: (a -> R b) -> R a -> R b
//
// This package also offers two versions for each of those functions. The
// regular versions `Fmap`, `Bind`, `Apply` take a single parameter and return
// a function. They are curried implementations.
// There are also corresponding uncurried implementations `EagerFmap`,
// `EagerBind`, and `EagerApply`. As the names suggest, other than their curried
// counterparts, they produce results immediately.
// For now, the eager versions are implemented in terms of their curried
// counterparts because the author somewhat suspects the 🪄 compiler 🪄 should be
// able to elide the extra func objects.
//
// TODO(kdungs): Check whether the compiler elides extra funcs for eager versions.
//
// Additionally, this package also implements
//
//	kleisli :: (a -> R b) -> (b -> R c) -> a -> R c
//	zipWith :: (a -> b -> c) -> R a -> R b -> R c
//
// This package also has a sister package `baresult` (bare result) which
// implements all the monadic goodness directly on `(T, error)` without
// introducing a wrapper type.
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
// In the context of R as a monad, this is equivalent to `return` or `pure` in
// Haskell: return :: a -> m a.
func Of[T any](v T) R[T] {
	return Wrap(v, nil)
}

// OfErr[T] constructs an `R[T]` holding an error.
// Since `T` cannot be inferred from the error alone, calling code is required
// to provide it, e.g. `OfErr[int](errors.New("not a positive number"))`.
func OfErr[T any](err error) R[T] {
	return Wrap(*new(T), err)
}
