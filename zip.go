package result

// ZipWith combines two results into one by applying a binary function to both
// values in case they exist. Otherwise whichever error is encountered first.
// zipWith :: (a -> b -> c) -> R a -> R b -> R c
func ZipWith[A any, B any, C any](f func(A, B) C) func(R[A]) func(R[B]) R[C] {
	return func(r R[A]) func(R[B]) R[C] {
		a, err := r.Unwrap()
		if err != nil {
			return func(_ R[B]) R[C] {
				return OfErr[C](err)
			}
		}
		return func(r R[B]) R[C] {
			b, err := r.Unwrap()
			if err != nil {
				return OfErr[C](err)
			}
			return Of(f(a, b))
		}
	}
}
