// Package then implements composition around "result functions" of type `func[A, B any](A) (B, error)`.
//
// The general idea is to lazily chain together computations without having to write `if err != nil` a lot.
//
// The major downside of this approach is that handling `defer`ed cleanup is not well supported.
package then

type (
	// FN is a result-producing function, the basic building block of this library.
	FN[A, B any] func(A) (B, error)

	/* There are only three other kinds of functions with which we want to compose.*/

	// F0 is a function that consumes a value of type A without returning an error.
	// This would usually by a function that modifies its argument or produces an error-free side effect (e.g. `fmt.Printf` but not `fmt.Fprintf`).
	F0[A any] func(A)

	// FE is a function that consumes a value of type A and returns an error (or nil).
	// This would typically be a function that produces a side effect, e.g. writing to a file via `fmt.Fprintf`.
	FE[A any] func(A) error

	// F is an regular function from A to B.
	F[A, B any] func(A) B
)

/*
With the above definitions, these are the scenarios for composition:
 FN + FN = FN => Chain
 FN + F  = FN => Map
 FN + FE = FE => Do
 FN + F0 = FE => Do0

On top of that, there are also zip-likes. Here, `x2` denotes a binary function
of the same kind as its unary counterpart.
 FN2(FN, FN) = FN2 => Zip
 FE2(FN, FN) = FE2 => Merge

Due to the lack of support for variadic (heterogenous) type arguments in Go, we
cannot implement the generic case for arbitrary higher arities. There are two
ways calling code can deal with higher arities:
 1. Wrapping parameters and return types in `struct`s.
 2. Folding through repeated application of `Zip` / `Merge`.
*/
// TODO(kdungs): Add examples for higher arities.

// Chain composes two result functions into a new result function.
func Chain[A, B, C any](f FN[A, B], g FN[B, C]) FN[A, C] {
	return func(a A) (C, error) {
		b, err := f(a)
		if err != nil {
			return *new(C), err
		}
		return g(b)
	}
}

// Lift takes a regular (error-free) function and elevates it to a result function.
// The returned error is always nil.
func Lift[A, B any](f F[A, B]) FN[A, B] {
	return func(a A) (B, error) {
		return f(a), nil
	}
}

// Map combines a result function with an error-free one.
// This is a convenience wrapper around `Chain` + `Lift`.
func Map[A, B, C any](f FN[A, B], g F[B, C]) FN[A, C] {
	return Chain(f, Lift(g))
}

// Do combines a result function with a consuming function that returns an error.
func Do[A, B any](f FN[A, B], g FE[B]) FE[A] {
	return func(a A) error {
		b, err := f(a)
		if err != nil {
			return err
		}
		return g(b)
	}
}

// Lift takes a consuming function that doesn't return an error and elevates it to a consuming function that does.
// The returned error is always nil.
func Lift0[A any](f F0[A]) FE[A] {
	return func(a A) error {
		f(a)
		return nil
	}
}

// Map combines a result function with a consuming function that doesn't return an error.
// This is a convenience wrapper around `Do` + `Lift0`.
func Do0[A, B any](f FN[A, B], g F0[B]) FE[A] {
	return Do(f, Lift0(g))
}

// Zip combines two result functions by applying a binary result function to their (non-error) results.
// If one of the results is an error, that error is returned instead.
func Zip[A, B, C, D, E any](f FN[A, B], g FN[C, D], with func(B, D) (E, error)) func(A, C) (E, error) {
	return func(a A, c C) (E, error) {
		b, err := f(a)
		if err != nil {
			return *new(E), err
		}
		d, err := g(c)
		if err != nil {
			return *new(E), err
		}
		return with(b, d)
	}
}

// Merge combines two result functions by applying a binary consuming function that returns an error to their (non-error) results.
// If one of the results is an error, that error is returned instead.
func Merge[A, B, C, D any](f FN[A, B], g FN[C, D], with func(B, D) error) func(A, C) error {
	return func(a A, c C) error {
		b, err := f(a)
		if err != nil {
			return err
		}
		d, err := g(c)
		if err != nil {
			return err
		}
		return with(b, d)
	}
}
