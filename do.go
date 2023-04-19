package result

// Do applies `f` to the value contained in `r` if present or otherwise returns
// the contained error.
func Do[T any](r R[T], f func(T)) error {
	rerr, ok := r.(rErr[T])
	if ok {
		return rerr.err
	}
	v, _ := r.Unwrap()
	f(v)
	return nil
}

// DoE applies `f` to the value contained in `r` if present and returns its
// result. Otherwise returns the error contained in `r`.
func DoE[T any](r R[T], f func(T) error) error {
	rerr, ok := r.(rErr[T])
	if ok {
		return rerr.err
	}
	v, _ := r.Unwrap()
	return f(v)
}

// DoZip applies `f` to the values contained in `ra` and `rb` if both are
// present. Otherwise returns the first error encountered (`ra` then `rb`).
func DoZip[A, B any](ra R[A], rb R[B], f func(A, B)) error {
	aerr, ok := ra.(rErr[A])
	if ok {
		return aerr.err
	}
	a, _ := ra.Unwrap()
	berr, ok := rb.(rErr[B])
	if ok {
		return berr.err
	}
	b, _ := rb.Unwrap()
	f(a, b)
	return nil
}

// DoZipE applies `f` to the values contained in `ra` and `rb` if both are
// present and returns its result. Otherwise returns the first error
// encountered (`ra` then `rb`).
func DoZipE[A, B any](ra R[A], rb R[B], f func(A, B) error) error {
	aerr, ok := ra.(rErr[A])
	if ok {
		return aerr.err
	}
	a, _ := ra.Unwrap()
	berr, ok := rb.(rErr[B])
	if ok {
		return berr.err
	}
	b, _ := rb.Unwrap()
	return f(a, b)
}
