package baresult

// Bind.
// TODO(kdungs): Documentation for `Bind`.
// bind :: (a -> R b) -> R a -> R b
func Bind[A any, B any](f func(A) (B, error)) func(A, error) (B, error) {
	return func(a A, err error) (B, error) {
		if err != nil {
			return *new(B), err
		}
		return f(a)
	}
}

// EagerBind is the uncurried, eager version of `Bind`.
// This exist because it (can be optimized to) avoid(s) an extra func.
// TODO(kdungs): Check whether the compiler actually elides the extra func.
// This is a bit useless since we cannot call it direcly with the result of a function...
func EagerBind[A any, B any](f func(A) (B, error), a A, err error) (B, error) {
	return Bind(f)(a, err)
}
