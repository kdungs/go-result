# Result types for Go

_Because `(T, error)` can be considered a ✨ monad ✨._

> Your scientists were so preoccupied with whether or not they could that they didn't stop to think if they should!

Imagine Go code without `if err != nil`. Would this be sacrilege? You might think so. For me, it was a fun exercise to get to know Go generics. Don't try this ~~at home~~ in prod.

The central idea is to treat `(T, error)` as a monad. There are two separate packages that implement this idea in different ways:

1.  [`package result`](result/) exposes a dedicated type `R[T]` that wraps `(T, error)`. On top of that, it implements functions to perform computations on `R[T]`. Since Go does not have generic member functions, those are free functions.
2.  [`package then`](then/) offers functionality to (lazily) compose functions that return `(T, error)`.

There's actually an [ongoing discussion around how to improve error handling in Go 2](https://go.googlesource.com/proposal/+/master/design/go2draft-error-handling.md).