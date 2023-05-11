# result

Why would this be useful? Take a look at [example_test.go](example_test.go).

## Caveats

 - No shortcuts! If an error happens early on, all computations still happen in
   the form of forwarding the error. With the current implementation this also
   means constructing a few default values on the way.
 - `defer Cleanup()` becomes a bit clunky, e.g.

```go
fh := result.Wrap(os.Open(fname))
defer result.Do(fh, func(f *os.File) { f.Close() })
```
