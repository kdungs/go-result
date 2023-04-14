package baresult

// Kleisli composition for functions that return `(T, error)`.
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
