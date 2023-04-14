package baresult

// Kleisli composition for functions that return `R[T]`.
// kleisli :: (a -> R b) -> (b -> R c) -> a -> R c
func Kleisli[A any, B any, C any](
	f func(A) (B, error),
	g func(B) (C, error),
) func(A) (C, error) {
	boundG := Bind(g)
	return func(a A) (C, error) {
		return boundG(f(a))
	}
}

// EagerKleisli is the uncurried, eager version of `Kleisli`.
// This exist because it (can be optimized to) avoid(s) an extra func.
// TODO(kdungs): Check whether the compiler actually elides the extra func.
func EagerKleisli[A any, B any, C any](
	f func(A) (B, error),
	g func(B) (C, error),
	a A,
) (C, error) {
	return Bind(g)(f(a))
}
