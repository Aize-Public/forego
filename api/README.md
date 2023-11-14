# `api`

Simplifies API design, documentation and testing by binding a go `type` to an API using `tags`.

```go
type WordFilter struct {
	Blacklist *regexp.Regexp

	R     api.Request `url:"/api/wordfilter/v1"`
	In    string      `api:"in,required" json:"in"`
	Out   string      `api:"out" json:"out"`
	Count int         `api:"out" json:"count"`
}

func (this *WordFilter) Do(c ctx.C) error {
	this.Out = this.Blacklist.ReplaceAllStringFunc(this.In, func(bad string) string {
		this.Count++
		return "***"
	})
	return nil
}
```

The above object can easily be "unit" tested by itself:

```go
func TestWordFilter(t *testing.T) {
	re := regexp.MustCompile(`(bad|worse|worst)`)
	out := api.Test(t, &WordFilter{
		Blacklist: re,
		In:        "ok, bad or worse",
	})
	test.EqualsStr(t, "ok, *** or ***", out.Out)
	test.EqualsGo(t, 2, out.Count)
}
```

And can be used with `forego/http` to directly
generate an API path:

```go
	s := http.NewServer(c)
	_ = s.RegisterAPI(c, &WordFilter{
		Blacklist: regexp.MustCompile(`(foo|bar)`), // this will be copied by ref for each request
	})
```

Which exposes the above object to `/api/wordfilter/v1` and also generate the OpenAPI documentation for it [see `/http` for more details].

## Handler

While normally you would use `api` directly from `http.Server` as shown above, you can also use it directly. Example for the server side:

```go
  // this is best to be done on startup…
	h, err := api.NewHandler(c, &WordFilter{})
	ser = h.Server() // can be used to perform Marshal/Unmarshal from the server side

  // then on each request…
	onRequest := func(c ctx.C, req api.ServerRequest, res api.ServerResponse) error {

		// un marshal the request into a new *WordFilter
		op, err := ser.Recv(c, req)
		if err != nil {
			return err
		}

		// call *WordFilter.Do()
		err = op.Do(c)
		if err != nil {
			return err
		}

		// marshal back the response
		return ser.Send(c, op, res)
	}
```

## Authentication

When `h.Server().Recv()`, the method `api.ServerRequest.Auth()` is called, which can then return an error.


## State

As in the example above, 


## `api.JSON{}`

This object is an helper for building `api.ServerRequest`, `api.ServerReponse` and the equivalent for the client.

It can be used on either ends to help building clients, unless you use `http.Server.Register` directly and avoid all the boilerplate.


## Field tags

Beside the common `json` tag used normally to marshal JSON data, there are few new tags you need to care about:

### `api`

This tag could be either `in`, `in,required`, `out`, `both,required` or by default `both`

When unmarshalling a request on a server, only `in`, `in,required` or `both` will be processed.

Moreover, if `in,required`, or `both,required` but the field is zero, a `400` error is returned

Only fields tagged as `out` or `both` will be marshalled back to the client.

Conversely, the same happen on the client side, just with the roles reversed.


### `auth`

If a field is marked as `auth`, it will be unmarshalled used the UID provided by the server side plumbing.

Moreover, if `auth,required` a 403 should be returned if the authentication token is missing


### `url`

Only one field should have the tag `url`, which provides a definition for the entrypoint of the API, and if fully qualified, can be used from clients 
directly without overrides.



