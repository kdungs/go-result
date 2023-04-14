package result

// Fmap.
// TODO(kdungs): Documentation for `Fmap`.
// fmap :: (a -> b) -> R a -> R b
func Fmap[A any, B any](f func(A) B) func(R[A]) R[B] {
	return func(r R[A]) R[B] {
		v, err := r.Unwrap()
		if err != nil {
			return OfErr[B](err)
		}
		return Of(f(v))
	}
}
