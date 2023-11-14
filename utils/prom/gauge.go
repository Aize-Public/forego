package prom

import (
	"fmt"
	"io"
	"sync/atomic"

	"github.com/Aize-Public/forego/utils/sync"
)

type gaugeEntry interface {
	value() float64
}

// wrap func() V as gaugeEntry[V]
type gaugeFunc func() float64

var _ gaugeEntry = (gaugeFunc)(nil)

func (this gaugeFunc) value() float64 { return this() }

/* prom.Gauge{} allows to set a function that handle a specific label */
type Gauge struct {
	Desc   string
	Labels []string
	val    sync.Map[string, gaugeEntry]
}

// create a int64 counter to be used as a gauge, will panic if not int64 or an incompatible guage was used before
func (this *Gauge) Counter(labels ...string) *GaugeCounter {
	var x gaugeEntry = &GaugeCounter{}
	// NOTE this will panic if a previous Gauge is there which is not a GaugeCounter
	return any(this.val.GetOrStore(stringify(this.Labels, labels), x)).(*GaugeCounter)
}

type GaugeCounter struct {
	val int64
}

var _ gaugeEntry = &GaugeCounter{}

func (this *GaugeCounter) value() float64 {
	return float64(atomic.LoadInt64(&this.val))
}

func (this *GaugeCounter) Inc(amt int) *GaugeCounter {
	atomic.AddInt64(&this.val, int64(amt))
	return this
}

func (this *GaugeCounter) Dec(amt int) *GaugeCounter {
	atomic.AddInt64(&this.val, int64(-amt))
	return this
}

func (this *Gauge) SetFunc(f func() float64, labels ...string) *Gauge {
	if f == nil {
		this.val.Delete(stringify(this.Labels, labels))
	} else {
		this.val.Store(stringify(this.Labels, labels), gaugeFunc(f))
	}
	return this
}

func (this *Gauge) Print(name string, w io.Writer) error {
	first := true
	return this.val.RangeErr(func(l string, x gaugeEntry) error {
		if first {
			first = false
			_, err := fmt.Fprintf(w, "# HELP %s %s\n", name, this.Desc)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(w, "# TYPE %s gauge\n", name)
			if err != nil {
				return err
			}
		}
		if l == "" {
			_, err := fmt.Fprintf(w, "%s %v\n", name, x.value())
			return err
		} else {
			_, err := fmt.Fprintf(w, "%s{%s} %v\n", name, l, x.value())
			return err
		}
	})
}
