# `ctx`

We expand from `context.Context` with few features and quality of life:
* mostly we assume we always want a `ctx.C` everywhere: because tracing, span, tags, loggers, etc
* contexts are not just for cancel, they can store environment and even configurations

## `c ctx.C` instead of `ctx context.Context`

Just cosmetic, I find it better to keep the context to the minimum, and used everywhere. And we never really used `c` for connections, clients or channels anyways


## Tags and `ctx/log`

Each context has a bag of tags, which can added along the way. Those will be printed in each log messages, which make it particularly useful for
things like `CorrelationID`, `auth` or any context which will help debugging from a log message.

It also make it coherent when using other libraries, since they will still carry over the context.

All logging is `JSON` lines, e.g.:

```
{"level":"debug","src":"github.com/Aize-Public/forego/http/server.go:83","time":"2023-06-01T07:18:31.007411033+02:00","message":"listening to :8080","tags":{"service":"viewer"}}
```

May be wise to use a log viewer like `https://github.com/ohait/jl`   


## `ctx.Err`

```go
  return ctx.NewErrorf(c, "my error wrapping %w", err)
```

Having a wrapping error that provide a stack trace has proven formidable when debugging or operating.

when the logger find a `ctx.Err` as an argument (or anything wrapping it) it will print the stack trace as part of the error message.


## Caveats

Generating stack traces is expensive in go, so don't use wrapping errors if you expect to ignore them often.
