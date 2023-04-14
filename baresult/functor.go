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

// EagerFmap is the uncurried, eager version of `Fmap`.
// This exist because it (can be optimized to) avoid(s) an extra func.
// TODO(kdungs): Check whether the compiler actually elides the extra func.
// This is a bit useless since we cannot call it direcly with the result of a function...
func EagerFmap[A any, B any](f func(A) B, a A, err error) (B, error) {
	return Fmap(f)(a, err)
}
