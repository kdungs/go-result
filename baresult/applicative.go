package baresult

// Apply.
// TODO(kdungs): Documentation for `Apply`.
// apply :: R (a -> b) -> R a -> R b
// TODO(kdungs): Reconsider implementation of `Apply` as it's not truly lazy.
func Apply[A any, B any](f func(A) B, err error) func(A, error) (B, error) {
	if err != nil {
		return func(_ A, _ error) (B, error) {
			return *new(B), err
		}
	}
	return func(a A, err error) (B, error) {
		if err != nil {
			return *new(B), err
		}
		return f(a), nil
	}
}
