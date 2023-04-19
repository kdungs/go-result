package result

// Map applies an error-free function to a result of type `R[A]`. In case the
// `R[A]` is holding a value, the returned result will contain the value
// returned from the function call. In case the result is holding an error, the
// returned result will be holding the same error.
// In functional programming terms, this is the same as `fmap` over a functor.
func Map[A, B any](r R[A], f func(A) B) R[B] {
	rerr, ok := r.(rErr[A])
	if ok {
		return OfErr[B](rerr.err)
	}
	a, _ := r.Unwrap()
	return Of(f(a))
}

// MapR applies a function whose return type is also a result to a result of type `R[A]`. In case the `R[A]` is holding a value, the returned result is the result of the function call. In case the result is holding an error, the returned result will be holding the same error.
// In functional programming terms, this corresponds to `bind` over a monad.
func MapR[A, B any](r R[A], f func(A) R[B]) R[B] {
	rerr, ok := r.(rErr[A])
	if ok {
		return OfErr[B](rerr.err)
	}
	a, _ := r.Unwrap()
	return f(a)
}

// MapE does the same as `MapR` but works for regular Go functions that return
// a value and an error.
func MapE[A, B any](r R[A], f func(A) (B, error)) R[B] {
	rerr, ok := r.(rErr[A])
	if ok {
		return OfErr[B](rerr.err)
	}
	a, _ := r.Unwrap()
	return Wrap(f(a))
}
