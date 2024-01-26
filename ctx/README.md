# `ctx`

We expand from `context.Context` with a few features and quality of life:
* mostly we assume we always want a `ctx.C` everywhere: because tracing, span, tags, loggers, etc
* contexts are not just for cancel, they can store environment, configuration and setting and overrides on each requests

## Why `c ctx.C` instead of `ctx context.Context`?

Just cosmetic, I find the traditional way distracting, especially when used everywhere.


## Tags and `ctx/log`

Each context has a bag of tags, which can added along the way. Those will by default be added to each log message,
which make them particularly useful for things like `CorrelationID`, `auth` or any context which will help debugging from a log message.

It also make it coherent when using other libraries, since they will still carry over the context.

By default, all logging is `JSONL`, e.g.:

```json
{"level":"debug","src":"github.com/Aize-Public/forego/http/server.go:83","time":"2023-06-01T07:18:31.007411033+02:00","message":"listening to :8080","tags":{"service":"viewer"}}
```

It may be wise to use a log viewer like `https://github.com/ohait/jl`


### `slog.Logger`

Alternatively, you can add your own `slog.Logger` to the context, and that one will be used instead (also in other forego libraries):

```go
myLogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
c = log.WithSlogLogger(c, myLogger)
log.Infof(c, "Hello world") // this will then be handled by myLogger, which in this example means it will be printed as a slog default JSON to stdout
```


### `LogFunc`

For even more control, you can add a custom `LogFunc` instead, which will bypass the automatic logging of tags.
Then you can also add a "Helper" function, and this example shows how that can be useful in tests:

```go
func TestSomething(t *testing.T) {
  c = log.WithHelper(c, t.Helper) // ensures that the correct src line will be logged by t.Logf
  c = log.WithLogFunc(c, func(c ctx.C, level slog.Level, src, f string, args ...any) {
    t.Helper()
    t.Logf("%s: %s", level, fmt.Sprintf(f, args...))
  })
  test.EqualsJSON(c, []any{1}, []int{1}) // the custom LogFunc is then passed along to the enc library via the context
}
```


## `ctx.Error`

```go
  return ctx.NewErrorf(c, "my error wrapping %w", err)
```

Having a wrapping error that provide a stack trace has proven formidable when debugging or operating.

When the logger find a `ctx.Error` as an argument (or anything wrapping it), it will add a stack trace to the "error" tag.


## Caveats

Generating stack traces is expensive in go, so don't use wrapping errors if you expect to ignore them often.

## TODO

finish setting up for opentelemetry (span) and tracking
