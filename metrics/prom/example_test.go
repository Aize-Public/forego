package prom

import (
	"bytes"
	"testing"
)

func TestBasic(t *testing.T) {
	m := Register(&Histogram{
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
	m.Observe(0.123, "foo", "bar")

	w := &bytes.Buffer{}
	//Handler()
	//err := m.Print(w)
	t.Logf("buf: %s", w.String())
}

func BenchmarkHistogram(b *testing.B) {
	m := &Histogram{
		Name:   "foo_bar",
		Desc:   "Foo Bar example",
		Labels: []string{"op", "loc"},
		Buckets: []float64{
			0.001, 0.002, 0.005,
			0.01, 0.02, 0.05,
			0.1, 0.2, 0.5,
			1, 2, 5,
		},
	}
	ops := []string{"read", "write"}
	loc := []string{"foo", "bar", "cuz"}
	b.Logf("bench: %d with %v, %v", b.N, ops, loc)
	for i := 0; i < b.N; i++ {
		m.Observe(0.123, ops[i%len(ops)], loc[i%len(loc)])
	}
}
