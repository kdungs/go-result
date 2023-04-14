package baresult

// Fmap.
// TODO(kdungs): Documentation for `Fmap`.
// fmap :: (a -> b) -> R a -> R b
func Fmap[A any, B any](f func(A) B) func(A, error) (B, error) {
	return func(a A, err error) (B, error) {
		if err != nil {
			return *new(B), err
		}
		return f(a), nil
	}
}
