package prom

import (
	"fmt"
	"io"

	"github.com/Aize-Public/forego/sync"
)

type Val interface{ int | int64 | float64 }

/* prom.Gauge{} allows to set a function that handle a specific label */
type Gauge[V Val] struct {
	Name   string
	Desc   string
	Labels []string
	val    sync.Map[string, func() V]
}

func (this *Gauge[V]) SetFunc(f func() V, labels ...string) *Gauge[V] {
	if f == nil {
		this.val.Delete(stringify(this.Labels, labels))
	} else {
		this.val.Store(stringify(this.Labels, labels), f)
	}
	return this
}

func (this *Gauge[V]) Print(w io.Writer) error {
	first := true
	return this.val.Range(func(l string, f func() V) error {
		if first {
			first = false
			_, err := fmt.Fprintf(w, "# HELP %s %s\n", this.Name, this.Desc)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(w, "# TYPE %s gauge\n", this.Name)
			if err != nil {
				return err
			}
		}
		if l == "" {
			_, err := fmt.Fprintf(w, "%s %v\n", this.Name, f()) // NOTE(oha) f() could be called concurrently, but i see no reason to protect from that
			return err
		} else {
			_, err := fmt.Fprintf(w, "%s{%s} %v\n", this.Name, l, f())
			return err
		}
	})
}
