package prom_test

import (
	"bytes"
	"testing"

	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/prom"
)

func TestGauge(t *testing.T) {
	m := prom.Gauge{
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
	m.Print(t.Name(), w)
	t.Logf("full: \n%s", w.String())
	test.Contains(t, w.String(), "/foo")
	test.NotContains(t, w.String(), "write")
	test.Contains(t, w.String(), "0.3")
}

func TestGaugeCounter(t *testing.T) {
	m := prom.Gauge{
		Labels: []string{"path", "op"},
	}
	c1 := m.Counter("/foo", "one")
	c1.Inc(42)
	{
		w := &bytes.Buffer{}
		m.Print(t.Name(), w)
		test.Contains(t, w.String(), "one")
		test.Contains(t, w.String(), "42")
	}
	c1.Dec(17)
	{
		w := &bytes.Buffer{}
		m.Print(t.Name(), w)
		test.Contains(t, w.String(), "25")
	}
}
