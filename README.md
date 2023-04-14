# go-result

_Result types for Go; because `(T, error)` can be considered a ✨ monad ✨._

> Your scientists were so preoccupied with whether or not they could that they didn't stop to think if they should!

**Disclaimer**: This whole thing is purely for (mostly my personal) education. Don't try this ~~at home~~ in prod.

## Background 
So... I was playing around with Rust. I'm not gonna lie: It's glorious. It reminded me of the good old Haskell days. Back then, I was inspired to implement a [check monad (aka Result, StatusOr, Either) for Python](https://github.com/kdungs/python-mcheck) and [Kleisli composition in C++](https://github.com/kdungs/cpp-kleisli-composition/). I thought to myself: Why not do something like this for Go? Purely as an experiment, of course. No sane Gopher would ever want to rid their code of all instances of

```go
if err := f(); err != nil {
    return nil, err
}
```

## Design
I came up with a number of different approaches that are discussed below under [Alternatives considered](#alternatives-considered). Ultimately, I decided to implement two of them. Both are based on the idea that `(T, error)` is equivalent to [Haskell's `Either Error a`](https://wiki.haskell.org/Handling_errors_in_Haskell#Error_using_the_Either_type) or [Rust's `Result<T>`](https://doc.rust-lang.org/std/result/). The `result` package exposes a dedicated generic type `R[T]` that wraps `(T, error)` and a number of free functions to make `R[T]` a ✨ monad ✨ (and an applicative functor). The `baresult` (bare result) package treats `(T, error)` as the monad, directly.

 * `Fmap` which corresponds to [`fmap :: Functor f => (a -> b) -> f a -> f b`](). Note that this the same as "lifting 🏋️ a function into the monad", see also [`liftM`](https://hackage.haskell.org/package/base-4.18.0.0/docs/Control-Monad.html#v:liftM).
  * `Apply` which corresponds to [`(<*>) :: Applicative f => f (a -> b) -> f a -> f b`](https://hackage.haskell.org/package/base-4.18.0.0/docs/Control-Applicative.html#v:-60--42--62-).
 * `Bind` which corresponds to [`(=<<) :: Monad m => (a -> m b) -> m a -> m b`](https://hackage.haskell.org/package/base-4.18.0.0/docs/Control-Monad.html#v:-61--60--60-) but **not** [`(>>=) :: m a -> (a -> m b) -> m b`](https://hackage.haskell.org/package/base-4.18.0.0/docs/Control-Monad.html#v:-62--62--61-). This is a design choice, see [discussion on currying below](#currying-and-go).

In addition, both packages also implement

 * `Kleisli` composition like [`(>=>) :: Monad m => (a -> m b) -> (b -> m c) -> a -> m c`](https://hackage.haskell.org/package/base-4.18.0.0/docs/Control-Monad.html#v:-62--61--62-)
 * `ZipWith` like `zipWith :: (a -> b -> c) -> R a -> R b -> R c`

For all of those functions, there's always the _default, lazy, curried_ and the _eager, uncurried_ version. The default versions all take a `func` and return another `func` which then acts on `R[T]`. This can be used to compose chains of functions and only actually execute them when a result is needed. E.g. `result.Fmap[A, B]` takes a `func(A) B` and returns a `func(result.R[A]) result.R[B]` while `result.EagerFmap[A, B]` takes both a `func(A) B` and a `result.R[A]` and immediately returns a `R[B]`.

### `result.R[T]`

 * `Of` which corresponds to [`return :: a -> m a`](https://hackage.haskell.org/package/base-4.18.0.0/docs/Control-Monad.html#v:return).
 * `OfErr` which does the same but for the error case. Since `T` cannot be inferred from the error alone, calls to `OfErr` have to specify `T` explicitly, e.g. `OfErr[float](err)`.
 * `Wrap` which takes `(T, error)` and returns an `R[T]`

### `baresult`

_TBD._

 * The `Eager` versions become a bit useless since we cannot e.g. do `EagerFmap(SomeFunc, SomeFuncThatReturnsValueAndError())` as the result of the second will not be unpacked. It does however work if the function call is the only parameter, e.g. `Fmap(SomeFunc)(SomeFuncThatReturnsValueAndError())`.

### Currying and Go
_TBD: Some discussion around why we return functions and why we flipped the order of arguments for `Bind`._

## Alternatives considered

### Getters for the value and error
The first draft contained `func (R[T])Err() error` and `func (R[T])Val() T` to allow direct access to the stored error and value. I found that this doesn't actually improve ergonomics since calling code would always have to call both of them anyway to check whether the `R[T]` is holding a value or an error. Also it requires a design decision around whether calling `Val()` while an error is present should cause a `panic` or just silently return the default value... `Unwrap` is more in line with standard Go and `Or` is a nice ergonomic option that _borrows_ from Rust's [`Result<T>::or`](https://doc.rust-lang.org/std/result/enum.Result.html#method.or).

### Storing `*T` instead
Instead of storing a value of type `T` inside `R[T]`, we could store a pointer. This would allow us to write `Val()` from the previous section without having to worry about whether to `panic` or not. In some cases, this might help avoid copies and improve performance (measurement and citation needed).

However, the ergonomics would suffer a lot. All functions dealing with `R[T]` would have to accept and return `*T`. Think about how clumsy something _primitive_ as `R[int]` would feel. With the implemented version, callers can decide to use a pointer by explicitly writing `R[*T]`.

We also don't need to worry about defaukt values because, fortunately, cheap default values are a paradigm in Go and the `*new(T)` 🦎 idiom is kind of neat (shoutout to Mat Ryer and the GoTime podcast).

### Interface `R[T]` with unexported implementations
Instead of a concrete type `R[T]`, we could define a generic interface

```go
type R[T any] interface {
    Unwrap() (T, error)
    Or(T) T
}
```

with two implementations for the error and value cases:

```go
type rError[T any] struct {
    err error
}

func (r rError[T]) Unwrap() (T, error) {
    return *new(T), r.err
}

func (r rError[T]) Or(d T) T {
    return d
}

type rValue[T any] struct {
    val T
}

func (r rValue[T]) Unwrap() (T, error) {
    return r.val, nil
}

func (r rValue[T]) Or(_ T) T {
    return r.val
}
```

The cool thing about this is that we'd never have to deal with default values (and pay for the cost of storing them) until an eventual call to `Unwrap` on an error.

Honestly, this might be the better design! Originally I had defined the interface in terms of `Err() error` and `Val() T` and it was clumsy around `panic`ing for `Val()` on the error-holding struct. Also exporting an interface and hiding concrete types is [kind of considered an anti-pattern in Go](https://github.com/golang/go/wiki/CodeReviewComments#interfaces) which is why I chose the `R[T] struct` instead. But it might be worth revisiting this...

## Future work

### n-ary wrapping and lifting 🏋️
I mean, do you even lift? More to the point: the current design has two shortcomings.

 1. `Wrap` wraps around the _result_ of a function call in the form of `(T, error)`. The function has to have already run, so the wrapping is eager. It would be great if we could have wrapping with lazy execution, instead. Imagine `WrapFunc[T any](f(...) (T, error)) func(...) R[T]`. The `...` will be discussed below.
 2. Similarly, we only have unary `Fmap` which lifts `func(A) B` into `func(R[A]) R[B]`. There's also `ZipWith` which can help dealing with binary functions but what about higher arities? Can we define a generic `Lift(func(Ts...) (T, error)) func(R[Ts...]) R[T]`? Think about [C++ parameter packs](https://en.cppreference.com/w/cpp/language/parameter_pack).

Unfortunately, Go's generics don't support what C++ calls [variadic templates or parameter packs](https://en.cppreference.com/w/cpp/language/parameter_pack). We cannot write the aforementioned generic `Lift` function because we have no way to express

 * a list of heterogenous types e.g. `Lift[Ts ...any]()`
 * operations on such a list, e.g. `R[T] for every T in Ts` as pack expansion would allow us to do in C++.

But, it's trivial to implement n-ary `WrapFunc` and `Lift` for a given `n`:

```go
func WrapFunc0[T any](f func() (T, error)) func() R[T] {
    return func() R[T] {
        return result.Wrap(f())
    }
}
// Lift0 would actually be exacly the same as WrapFunc0.
// Lift2 would actually be exactly the same as ZipWith.

func WrapFunc3[T any, A any, B any, C any](func (A, B, C) (T, error)) func(A, B, C) R[T] {
    return func(a A, b B, c C) R[T] {
        return f(a, b, c)
    }
}

func Lift3[T any, A any, B any, C any](func (A, B, C) (T, error)) func(R[A], R[B], R[C]) R[T] {
    return func(ra R[A], rb R[B], rc R[C]) R[T] {
        // Or more elegantly in terms of Apply, Fmap, Lift2, whatever...
        a, err := ra.Unwrap()
        if err != nil {
            return result.OfErr[T](err)
        }
        b, err := ra.Unwrap()
        if err != nil {
            return result.OfErr[T](err)
        }
        c, err := ra.Unwrap()
        if err != nil {
            return result.OfErr[T](err)
        }
        return result.Wrap(f(a, b, c))
    }
} // Or even inductively... TBD.
```

So we could just manually (or through code generation) implement all versions of `WrapFunc` and `Lift` up to a specific `N`. While this might sound a bit ridiculous in 2023, this is actually how people used to do it in C++03, e.g. in Boost. If you don't believe me, take a look at ["C++ Template Metaprogramming: Concepts, Tools, and Techniques from Boost and Beyond" by by David Abrahams and Aleksey Gurtovoy](https://www.oreilly.com/library/view/c-template-metaprogramming/0321227255/).

### Ergonomics
As I mentioned previously, this whole thing is kind of an academic exercise. This is also why I stuck to the Haskell names for the functions. For a real-world library, we wouldn't want to require users to know Haskell so a name like `Then` instead of `Bind` might be better suited.

Also, it would be nice to be able to call methods on `R[T]` instead of having to use chains of free functions. Function composition like that means we kind of have to reverse the order, e.g.

```go
Fmap(Last)(Fmap(Second)(Fmap(First)(r))).Unwrap()
```

Nicer approaches would be

```go
r.Fmap(First).Fmap(Second).Fmap(Last).Unwrap()
```

or even

```go
result.Chain(First, Second, Last).Unwrap()
```

Unfortunately, as awesome as Go generics are, they currently don't support generic member functions. This means, we don't have a way to implement e.g. `Fmap` on `R[T]`. A truly generic `Chain` would also require parameter packs unless we want to go the same route as with n-ary `Wrap` and `Lift`.