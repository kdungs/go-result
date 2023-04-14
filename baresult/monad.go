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
