package baresult

// ZipWith combines two results into one by applying a binary function to both
// values in case they exist. Otherwise whichever error is encountered first.
// zipWith :: (a -> b -> c) -> R a -> R b -> R c
func ZipWith[A any, B any, C any](f func(A, B) C) func(A, error, B, error) (C, error) {
	return func(a A, errA error, b B, errB error) (C, error) {
		if errA != nil {
			return *new(C), errA
		}
		if errB != nil {
			return *new(C), errB
		}
		return f(a, b), nil
	}
}

// EagerZipWith is the uncurried, eager version of `ZipWith`.
// This exist because it (can be optimized to) avoid(s) an extra func.
// TODO(kdungs): Check whether the compiler actually elides the extra func.
func EagerZipWith[A any, B any, C any](f func(A, B) C, a A, errA error, b B, errB error) (C, error) {
	return ZipWith(f)(a, errA, b, errB)
}
