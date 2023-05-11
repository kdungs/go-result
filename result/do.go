package result

// Do applies `f` to the value contained in `r` if present or otherwise returns
// the contained error.
func Do[T any](r R[T], f func(T)) error {
	v, err := r.Unwrap()
	if err != nil {
		return err
	}
	f(v)
	return nil
}

// DoE applies `f` to the value contained in `r` if present and returns its
// result. Otherwise returns the error contained in `r`.
func DoE[T any](r R[T], f func(T) error) error {
	v, err := r.Unwrap()
	if err != nil {
		return err
	}
	return f(v)
}

// DoZip applies `f` to the values contained in `ra` and `rb` if both are
// present. Otherwise returns the first error encountered (`ra` then `rb`).
func DoZip[A, B any](ra R[A], rb R[B], f func(A, B)) error {
	a, err := ra.Unwrap()
	if err != nil {
		return err
	}
	b, err := rb.Unwrap()
	if err != nil {
		return err
	}
	f(a, b)
	return nil
}

// DoZipE applies `f` to the values contained in `ra` and `rb` if both are
// present and returns its result. Otherwise returns the first error
// encountered (`ra` then `rb`).
func DoZipE[A, B any](ra R[A], rb R[B], f func(A, B) error) error {
	a, err := ra.Unwrap()
	if err != nil {
		return err
	}
	b, err := rb.Unwrap()
	if err != nil {
		return err
	}
	return f(a, b)
}
