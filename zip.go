package result

// Zip combines two results into one by applying a binary function that returns
// a value to both values in case they exist. Otherwise returns whichever error
// is encountered first.
func Zip[A, B, C any](ra R[A], rb R[B], f func(A, B) C) R[C] {
	a, err := ra.Unwrap()
	if err != nil {
		return OfErr[C](err)
	}
	b, err := rb.Unwrap()
	if err != nil {
		return OfErr[C](err)
	}
	return Of(f(a, b))
}

// ZipR combines two results into one by applying a binary function returning a
// result to both values in case they both exist. Otherwise returns whichever
// error is encountered first.
func ZipR[A, B, C any](ra R[A], rb R[B], f func(A, B) R[C]) R[C] {
	a, err := ra.Unwrap()
	if err != nil {
		return OfErr[C](err)
	}
	b, err := rb.Unwrap()
	if err != nil {
		return OfErr[C](err)
	}
	return f(a, b)
}

// ZipE does the same as `ZipR` but works for regular Go functions that return
// a value and an error.
func ZipE[A, B, C any](ra R[A], rb R[B], f func(A, B) (C, error)) R[C] {
	a, err := ra.Unwrap()
	if err != nil {
		return OfErr[C](err)
	}
	b, err := rb.Unwrap()
	if err != nil {
		return OfErr[C](err)
	}
	return Wrap(f(a, b))
}
