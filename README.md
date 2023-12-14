# forego

THIS MODULE IS STILL A WORK IN PROGRESS AND NOT READY FOR USE

Go framework to simplify testing, http/ws API, authentication/authorization, etc...

Some subfolders have a README with more details:

* [api](./api/)
* [test](./test/)
* [enc](./enc/)
* [http](./http/)
* [http/ws](./http/ws/)
* [shutdown](./shutdown/)
* [test](./test/)
* [utils/prom](./utils/prom/)


## `test`

Having a good testing library, means less to write, and better results.

All the tests will generate a log message when they succeed based on the arguments, most of the time this is enough to understand the test.

E.g. with this code:

```go
  err := foo(123)
  test.NoError(t, err)
```

The following log message is generated:   

```go
    my_test.go:123 ok: foo(123)
```

### `ast` parsing

Some functions like `test.NoError()` and `test.Assert()` will log the code which they have been invoked with. E.g.:

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
* no more test of bindings, just test the business logic
* OpenAPI automatically generated
* tight integration with http, WebSocket and other streaming libraries
* simple and no boiler plate

TODO

