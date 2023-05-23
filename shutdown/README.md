# Shutdown

Provides basic support for graceful shutdown.

The main goroutine will likely use `WaitForSignal()`:

```go
func main() {
  c, cf := ctx.Background()
  // your code here
  shutdown.WaitForSignal(c, cf)
}
```

While the other goroutines will likely keep the shutdown from completing until they are properly closed:

```go
func work(c ctx.C) {
  defer shutdown.Hold().Release() // obtain a hold immediately, and defer the release

  for {
    // main loop
    select {
    case <-c.Done():
      // abort

    case <-shutdown.Started():
      // graceful close and then

    case data, ok := <-io:
      if !ok {
        return // EOF
      }
      // do work
    }
  }
}
```

## `WaitForSignal()`

It blocks until there are no more active `Hold()`s, or 3 signals (INT, TERM or QUIT) has been detected.

When the first signal is detected calls `Begin()`

## `Begin()`

It begins the shutdown procedure. Anyone waiting for `<-shutdown.Started()` will be notified.

## `Started()`

Returns a channel which will be closed when the shutdown `Begin()`

## `Hold().Release()`

return a hold which prevents the shutdown to complete until it is `Released()`

## `Started5Sec()`

Similar to `Started()` but returns the channel returned will be closed 5 seconds later

