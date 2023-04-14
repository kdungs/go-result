package result

// Apply.
// TODO(kdungs): Documentation for `Apply`.
// apply :: R (a -> b) -> R a -> R b
// TODO(kdungs): Reconsider implementation of `Apply` as it's not truly lazy.
func Apply[A any, B any](r R[func(A) B]) func(R[A]) R[B] {
	f, err := r.Unwrap()
	if err != nil {
		return func(r R[A]) R[B] {
			return OfErr[B](err)
		}
	}
	return Fmap(f)
}
