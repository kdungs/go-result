package result

// Apply.
// TODO(kdungs): Documentation for `Apply`.
// apply :: R (a -> b) -> R a -> R b
func Apply[A any, B any](r R[func(A) B]) func(R[A]) R[B] {
	f, err := r.Unwrap()
	if err != nil {
		return func(r R[A]) R[B] {
			return OfErr[B](err)
		}
	}
	return Fmap(f)
}

// EagerApply is the uncurried, eager version of `Apply`.
// This exist because it (can be optimized to) avoid(s) an extra func.
// TODO(kdungs): Check whether the compiler actually elides the extra func.
func EagerApply[A any, B any](f R[func(A) B], a R[A]) R[B] {
	return Apply(f)(a)
}
