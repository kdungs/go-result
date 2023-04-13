package result

// Bind.
// TODO(kdungs): Documentation for `Bind`.
// bind :: (a -> R b) -> R a -> R b
func Bind[A any, B any](f func(A) R[B]) func(R[A]) R[B] {
	return func(r R[A]) R[B] {
		v, err := r.Unwrap()
		if err != nil {
			return OfErr[B](err)
		}
		return f(v)
	}
}

// EagerBind is the uncurried, eager version of `Bind`.
// This exist because it (can be optimized to) avoid(s) an extra func.
// TODO(kdungs): Check whether the compiler actually elides the extra func.
func EagerBind[A any, B any](f func(A) R[B], r R[A]) R[B] {
	return Bind(f)(r)
}
