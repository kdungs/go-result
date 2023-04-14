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
