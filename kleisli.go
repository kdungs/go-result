package result

// Kleisli composition for functions that return `R[T]`.
// kleisli :: (a -> R b) -> (b -> R c) -> a -> R c
func Kleisli[A any, B any, C any](f func(A) R[B], g func(B) R[C]) func(A) R[C] {
	boundG := Bind(g)
	return func(a A) R[C] {
		return boundG(f(a))
	}
}
