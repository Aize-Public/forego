# `prom` -- Prometheus lightweight replacement

Replaces `github.com/prometheus/client_golang` with a lighter version, no collector and simpler code

```go
var foobar = prom.Register(&prom.Histogram{
		Name:   "foo_bar",
		Desc:   "Foo Bar example",
		Labels: []string{"op", "loc"},
		Buckets: []float64{
			0.001, 0.002, 0.005,
			0.01, 0.02, 0.05,
			0.1, 0.2, 0.5,
			1, 2, 5,
		},
})

type FooBar struct {
  Op string
  Loc string
}

func (this FooBar) Observe(val float64) {
	  foobar.Observe(0.123, this.Op, this.Loc)
}
```

And in your http server:

```go
  s.HandleFunc("/metrics", prom.Handler())
```
