# `metrics`

Set of predefined metrics used globally by any service, e.g.:

```go
func fetch(c ctx.C, key string) (val string) {

	defer metrics.KeyValue{ // histogram
		Op:  "fetch",
		Src: util.Caller(1).FileLine(), // file.go:123
	}.ObserveSince(time.Now()) // time.Now() is executed immediately, but ObserveSince() is defer-ed -> we report how long fetch() took

  // do the fetching

  return "result"
}
```

## Prometheus

Prometheus protocol is very simple, it's based on old `statsd` and `graphite`. Unfortunately, the new features are based on
`protobuffers` and because of it, the default library imports a lot of other packages, some of which are quite old.

To maintain a small footprint, the Prometheus client was reverse engineers into `metrics/prom` and only the simplest necessary cases are supported.


