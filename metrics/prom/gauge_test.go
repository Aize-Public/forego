package prom_test

import (
	"bytes"
	"testing"

	"github.com/Aize-Public/forego/metrics/prom"
	"github.com/Aize-Public/forego/test"
)

func TestGauge(t *testing.T) {
	m := prom.Gauge[float64]{
		Name:   t.Name(),
		Labels: []string{"path", "op"},
	}
	m.SetFunc(func() float64 {
		return 3.14
	}, "/foo", "read")
	m.SetFunc(func() float64 {
		return 1.23
	}, "/foo", "write")
	m.SetFunc(func() float64 {
		return 0.3
	}, "/foo", "read")
	m.SetFunc(nil, "/foo", "write")

	w := &bytes.Buffer{}
	m.Print(w)
	t.Logf("full: \n%s", w.String())
	test.Contains(t, w.String(), "/foo")
	test.NotContains(t, w.String(), "write")
	test.Contains(t, w.String(), "0.3")
}
