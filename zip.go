package result

// Zip combines two results into one by applying a binary function that returns
// a value to both values in case they exist. Otherwise returns whichever error
// is encountered first.
func Zip[A, B, C any](ra R[A], rb R[B], f func(A, B) C) R[C] {
	aerr, ok := ra.(rErr[A])
	if ok {
		return OfErr[C](aerr.err)
	}
	a, _ := ra.Unwrap()
	berr, ok := rb.(rErr[B])
	if ok {
		return OfErr[C](berr.err)
	}
	b, _ := rb.Unwrap()
	return Of(f(a, b))
}

// ZipR combines two results into one by applying a binary function returning a
// result to both values in case they both exist. Otherwise returns whichever
// error is encountered first.
func ZipR[A, B, C any](ra R[A], rb R[B], f func(A, B) R[C]) R[C] {
	aerr, ok := ra.(rErr[A])
	if ok {
		return OfErr[C](aerr.err)
	}
	a, _ := ra.Unwrap()
	berr, ok := rb.(rErr[B])
	if ok {
		return OfErr[C](berr.err)
	}
	b, _ := rb.Unwrap()
	return f(a, b)
}

// ZipE does the same as `ZipR` but works for regular Go functions that return
// a value and an error.
func ZipE[A, B, C any](ra R[A], rb R[B], f func(A, B) (C, error)) R[C] {
	aerr, ok := ra.(rErr[A])
	if ok {
		return OfErr[C](aerr.err)
	}
	a, _ := ra.Unwrap()
	berr, ok := rb.(rErr[B])
	if ok {
		return OfErr[C](berr.err)
	}
	b, _ := rb.Unwrap()
	return Wrap(f(a, b))
}
