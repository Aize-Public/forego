# `http`

```go
  // create a new server
	s := http.NewServer(c, "example")

  // go built-in http/Handler
	s.Mux().HandleFunc("/test/one", func(w gohttp.ResponseWriter, r *gohttp.Request) {
		_, _ = w.Write([]byte(`"one"`))
	})

  // generic POST handler with error helpers
	s.HandleRequest("/test/two", func(c ctx.C, in []byte, r *gohttp.Request) ([]byte, error) {
		return enc.MarshalJSON(c, "one")
	})

  // using API library (see forego/api/)
	s.MustRegisterAPI(c, &MyAPI{})
  
  // listen to a port
	addr, err := s.Listen(c, "127.0.0.1:0")
  if err != nil {
    panic(err)
  }

  // wait for graceful shutdown
	shutdown.WaitForSignal(c, cf)
```

## `Server`

```go
	s := http.NewServer(c, "example")
```  

Creates a new server, it will internally setup several middleware which provide:
* logging (can be overridden with `s.OnRequest`
* context tags (`ua` for user agent, `path`, `http.addr` and `http.remote`)
* `ServeMux` (can be used with `s.Mux()`)
* `/live` which always return 204 (OK No Content)
* `/ready` which by default returns 204 (can be changed with `s.SetReady()`)
* `/openapi.json` serialized from `s.OpenAPI`

### `Mux()`

Returns the internal `ServeMux`, which can be then used to add new paths to the server using go built-in `http.Handler`

### `HandleRequest(pattern string, fn func(ctx.C, []byte, *http.Request) ([]byte, error))`

```go
	s.HandleRequest("/test/two", func(c ctx.C, in []byte, r *gohttp.Request) ([]byte, error) {
		return enc.MarshalJSON(c, "one")
	})
```

Helper which create a handler for the given function, which:
* read fully the request body
* call the given function
* if error is returned, the `trackingID` is returned with a 500
* if error is a http.Error, and the code is <500, then the error message is returned as well
* if no response, return 204
* otherwise set the response to `application/json` and send the data
* optionally gzip if more than 16KB and request accepts gzip

An entry is created in `s.OpenAPI` for the given path, and the `*openapi.PathInfo` is returned for further tweaking

### `Register(c, obj)` and `MustRegister(c, obj)`

```go
  err = s.Register(c, &MyApi)
```

Uses the [api](../api/) library to parse the given object, and expose the API.

It also update `s.OpenAPI` accordingly.


