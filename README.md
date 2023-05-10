# forego

Go framework to simplify testing, http/ws API, authentication/authorization, etc...

## `ctx`

We expand from `context.Context` with few features and quality of life:
* mostly we assume we always want a `ctx.C` everywhere: because tracing, span, tags, loggers, etc
* contexts are not just for cancel, they can store environment and even configurations

### `c ctx.C` instead of `ctx context.Context`

Just cosmetic, I find it better to keep the context to the minimum, and used everywhere. And we never really used `c` for connections, clients or channels anyways, do we?

### `ctx.Err`

Having a wrapping error that provide a stack trace has proven formidable when working with libraries. The logger expand any logged error accordingly. No longer guess from there the error come from, but still no extra log messages or long stack traces.

### Tags and `ctx/log`

Each context has a bag of tags, which can added along the way. Those will be printed in each log messages, which make it particularly useful for
things like `CorrelationID`, `auth` or any context which will help debugging from a log message.

It also make it coherent when using other libraries, since they will still carry over the context.

All logging is `JSON` lines

## `test`

Having a good testing library, means less to write, and better results.

all the tests will generate a log message when they succeed, most of the time this is enough to understand the test.

### `ast` parsing

Some functions like `test.NoError()` and `test.Assert()` will log the code which they have been invoked with. e.g.:

```go
  func TestX(t *testing.T) {
    everything := 42
    test.Assert(t, everything > 7*3)
  }
```

It will log a message like:
```
ok: everything > 7*3
```

In the case of error, it will print what was the function assigning to the error

## `api`

Framework to automatically create bindings and documentation for APIs:
* no more test of bindings, just test the logic
* OpenAPI automatically generated
* tight integration with http, ws and other streaming services
* simple and no boiler plate
