# `enc`

Replacement for `encoding/json`, providing an intermediate layer of abstraction between the encoded data and the typed data.


## Rationale

In the built in `encoding/json` library, everything gets converted between `[]byte` and specific types.

This means that when customization is needed, you end up implementing `json.Unmarshaler` or `json.RawMessage`.

This likely requires multiple scans to the data, and/or writing convoluted code.

Using an intermediate layer, which maps JSON to some intermediate types, allows for a single parsing of the `[]byte` data in.
After that, using and managing the intermediate types is easier, and less computational intensive.

Another way to look at it, the intermediate types can be used instead of `json.RawMessage`, making it easier to code and faster to execute.


## Marshal vs Encode

To simplify the documentation, we will use *marshal* when transforming an object into the interstitial types, and *encoding* when converting the
intermediate to a JSON `[]byte`.

Similarly, we say *decoding* when parsing the JSON `[]byte` and *unmarshalling* when coercing the interstitial into an object type

More details can be found in the distinct types.


## Types

### `enc.Node`

This is the generic interface, all interstitial types implements it and can be used as a replacement for `json.RawMessage`.

The following `struct` will only unmarshal a `enc.Map` or `enc.Pairs`:

```go
  // {
  //   "type": "my-type",
  //   "data": { … },
  //   "cmd": ["x", "y"],
  // }

  type Frame struct {
    Type string   `json:"type"` // Coerced into string
    Data enc.Node `json:"data"` // left as-is (enc.Map)
    Path []string `json:"cmd"`  // Unmarshalled further
  }
```

### `enc.Numeric` interface and `enc.Integer`, `enc.Float`

When decoding a `JSON number`, either a `enc.Integer` or a `enc.Float` will be returned, which both implements `enc.Numeric`.

They can unmarshal similarly to `encoding/json` when you use a `Decoder` with `UseNumber()`.

### `enc.String`

```go
  s := enc.String("foo")
```

### `enc.Bool`

```go
  b := enc.Bool(true)
```

### `enc.List`

```go
  l := enc.List{
    enc.String("answer"),
    enc.Integer(42),
  }
```

### `enc.Map`

Used for generic object types, the order of the field is not kept.

```go
  m := enc.Map{
    "type": enc.String("my-type"),
    "data": enc.Map{},
    "count": enc.Integer(42),
  }
  for k, v := range m {
    // type, data, count...
  }
```

### `enc.Pairs`

This is a special type that can Marshal itself as a `JSON object`, but is implemented as a list of pairs, which then guarantee the order.

To keep the usability of this library high, we opted to avoid Ordered-Maps which are clumsy to use, and instead allow you to choose between the 
fast `enc.Map`, or the ordered `enc.Pairs`.


### Custom `enc.Marshaler`, `enc.Unmarshaler` vs `json.Marshaler` and `json.Unmarshaler`

This library is compatible with `json.Marshaler` and `json.Unmarshaler`, but those interfaces requires to re-encode and re-decoded `[]byte`.

It is hence more efficient to use the new `enc.Marshaler` and `enc.Unmarshaler`.

Here an example object:

```go
type X struct {
	Type string
	Path []string
}

func (this X) String() string {
	s := this.Type
	sep := ":"
	for _, p := range this.Path {
		s += sep + url.QueryEscape(p)
		sep = "/"
	}
	return s
}

var xRE = regexp.MustCompile(`^([a-z]+):([a-z]+(?:\/[a-z]+)*)$`)

func (this *X) Parse(c ctx.C, s string) error {
	out := xRE.FindStringSubmatch(s)
	if len(out) == 0 {
		return ctx.NewErrorf(c, "invalid X: %q", s)
	}
	log.Warnf(c, "out: %#v", out)
	this.Type = out[1]
	this.Path = []string{}
	for _, p := range strings.Split(out[2], "/") {
		this.Path = append(this.Path, url.QueryEscape(p))
	}
	return nil
}
```

To make it use `String()` and `Parse()` when generating `enc.Node`:

```go
var _ enc.Marshaler = &X{}
var _ enc.Unmarshaler = &X{}

func (this X) MarshalNode(c ctx.C) (enc.Node, error) {
	return enc.String(this.String()), nil
}

func (this *X) UnmarshalNode(c ctx.C, n enc.Node) error {
	switch n := n.(type) {
	case enc.String:
		return this.Parse(c, string(n))
	default:
		return ctx.NewErrorf(c, "expected string, got %T", n)
	}
}
```

As you can see, it simply returns and `enc.String` to marshal, and only accept it back when unmarshalling.

To keep compatibility with `encoding/json` you might want to implement the relative versions of those methods too.

## `ctx.C`

One of the many reasons to implement and/or use this library is the integration with `context.Context` (or `ctx.C`)

Why this is so important can be debatable, but I often used `Context` as a way to inject specific logic into a generic framework.

It could be configuration at the top level, or request specific settings.

If those settings are needed in any of the `UnmarshalNode` or `MarshalNode`, to protect some data or create metrics, or automate subscriptions, you
can now use a pattern like the following:

```go
type ctxKey = struct{}

func SetupCallback(c ctx.C, fn func(x X) error) ctx.C {
  return ctx.WithValue(c, ctxKey{}, fn)
}

func NotifyCallback(c ctx.C, x X) error {
  fn := c.Value(c, ctxKey{})
  switch fn := fn.(type) {
  case nil:
    return nil
  case func(X) error:
    return fn(x)
  default:
    return ctx.NewErrorf(c, "unexpected %T: %v", fn, fn)
  }
}
```

Which can then be used on setup:
```go
  c = SetupCallback(c, func(x X) error {
    if x.Type == blacklist {
      return ctx.NewErrorf(c, "invalid type: %q", x.Type)
    }
    return nil
  })
```

Where blacklist might depends on the language of the user, for example.

It can then be triggered while unmarshalling the requests:

```go
func (this *X) UnmarshalNode(c ctx.C, n enc.Node) error {
  switch n := n.(type) {
  case enc.String():
    err := this.Parse(n)  
    if err != nil {
      return err
    }
    return NotifyCallback(c, *this) // if there is a callback, and returns an errors, unmarshalling is aborted
  default:
    return ctx.NewErrorf(c, "expected string, got %T", n)
  }
}
```

which means that now any API you have will fail if they contains a type `X` which has a `.Type` which is blacklisted, and
that blacklist might change per request, based on the user settings or permissions or so.
