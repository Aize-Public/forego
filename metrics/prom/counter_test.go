package prom_test

import (
	"bytes"
	"testing"

	"github.com/Aize-Public/forego/metrics/prom"
	"github.com/Aize-Public/forego/test"
)

func TestCounter(t *testing.T) {

	m := prom.Counter{
		Name:   t.Name(),
		Labels: []string{"path", "op"},
	}
	m.Observe(3.14, "/foo", "read")
	m.Observe(0.5, "/foo", "write")
	m.Observe(1.0, "/foo", "write")
	m.Observe(2, "/bar", "read")

	w := &bytes.Buffer{}
	m.Print(w)
	t.Logf("full: \n%s", w.String())
	test.Contains(t, w.String(), "/foo")
	test.Contains(t, w.String(), "1.5") // 0.5 + 1.0

}
