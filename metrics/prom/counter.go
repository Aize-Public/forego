package prom

import (
	"fmt"
	"io"

	"github.com/Aize-Public/forego/sync"
)

type Counter struct {
	Name   string
	Desc   string
	Labels []string
	val    sync.Map[string, *counter]
	Gauge  bool // set to true if this counter needs to be reported as a gauge (but see gauge.go for a better API)
}

type counter struct {
	m   sync.Mutex
	sum float64
}

func (this *Counter) Observe(val float64, labels ...string) {
	this.val.GetOrStore(stringify(this.Labels, labels), &counter{}).observe(val)
}

func (this *counter) observe(val float64) {
	//defer log.Printf("(%p).Observe(%v)", this, val)
	this.m.Lock()
	defer this.m.Unlock()
	this.sum += val
}

func (this *Counter) Print(w io.Writer) error {
	first := true
	return this.val.Range(func(l string, v *counter) error {
		if first {
			first = false
			_, err := fmt.Fprintf(w, "# HELP %s %s\n", this.Name, this.Desc)
			if err != nil {
				return err
			}
			if this.Gauge {
				_, err = fmt.Fprintf(w, "# TYPE %s gauge\n", this.Name)
				if err != nil {
					return err
				}
			} else {
				_, err = fmt.Fprintf(w, "# TYPE %s counter\n", this.Name)
				if err != nil {
					return err
				}
			}
		}
		if l == "" {
			_, err := fmt.Fprintf(w, "%s %f\n", this.Name, v.sum)
			return err
		} else {
			_, err := fmt.Fprintf(w, "%s{%s} %f\n", this.Name, l, v.sum)
			return err
		}
	})
}
