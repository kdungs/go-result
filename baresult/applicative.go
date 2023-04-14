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

// EagerApply is the uncurried, eager version of `Apply`.
// This exist because it (can be optimized to) avoid(s) an extra func.
// TODO(kdungs): Check whether the compiler actually elides the extra func.
// This is a bit useless since we cannot call it direcly with the result of a function...
func EagerApply[A any, B any](f func(A) B, errF error, a A, errA error) (B, error) {
	return Apply(f, errF)(a, errA)
}
