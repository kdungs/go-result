package baresult

// ZipWith combines two results into one by applying a binary function to both
// values in case they exist. Otherwise whichever error is encountered first.
// zipWith :: (a -> b -> c) -> R a -> R b -> R c
func ZipWith[A any, B any, C any](f func(A, B) C) func(A, error) func(B, error) (C, error) {
	return func(a A, err error) func(B, error) (C, error) {
		if err != nil {
			return func(_ B, _ error) (C, error) {
				return *new(C), err
			}
		}
		return func(b B, err error) (C, error) {
			if err != nil {
				return *new(C), err
			}
			return f(a, b), nil
		}
	}
}
