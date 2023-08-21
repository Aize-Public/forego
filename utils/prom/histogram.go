package prom

import (
	"fmt"
	"io"
	"time"

	"github.com/Aize-Public/forego/sync"
)

type Histogram struct {
	//Name    string
	Desc    string
	Buckets []float64
	Labels  []string
	val     sync.Map[string, *histogramTS]
}

var DefaultBuckets = []float64{.001, .0025, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30, 60, 60 * 5, 60 * 30} // 1ms .. 30m

type histogramTS struct {
	m       sync.Mutex
	buckets []int
	sum     float64
	count   int
}

func (this *Histogram) ObserveSince(start time.Time, labels ...string) {
	this.val.GetOrStore(stringify(this.Labels, labels), &histogramTS{}).
		observe(this.Buckets, time.Since(start).Seconds())
}

func (this *Histogram) Observe(val float64, labels ...string) {
	this.val.GetOrStore(stringify(this.Labels, labels), &histogramTS{}).
		observe(this.Buckets, val)
}

func (this *histogramTS) observe(le []float64, val float64) {
	//defer log.Printf("(%p).Observe(%v)", this, val)
	this.m.Lock()
	defer this.m.Unlock()

	this.count++
	this.sum += val

	for len(this.buckets) < len(le) {
		this.buckets = append(this.buckets, 0)
	}
	for i, le := range le {
		if val <= le {
			this.buckets[i]++
		}
	}
}

func (this *Histogram) Print(name string, w io.Writer) (err error) {
	first := true
	return this.val.Range(func(l string, v *histogramTS) error {
		if first {
			first = false
			_, err := fmt.Fprintf(w, "# HELP %s %s\n", name, this.Desc)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintf(w, "# TYPE %s histogram\n", name)
			if err != nil {
				return err
			}
		}
		return v.print(w, this.Buckets, name, l)
	})
}

func (this *histogramTS) print(w io.Writer, le []float64, name, labels string) (err error) {
	this.m.Lock()
	defer this.m.Unlock()
	if labels == "" {
		for i, le := range le {
			b := this.buckets[i]
			_, err = fmt.Fprintf(w, "%s_bucket{le=\"%v\"} %d\n", name, le, b)
			if err != nil {
				return err
			}
		}
		// +Inf is reported, even tho it's dup of "_sum"...
		_, err = fmt.Fprintf(w, "%s_bucket{le=\"+Inf\"} %d\n", name, this.count)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "%s_sum %f\n", name, this.sum)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%s_count %d\n", name, this.count)
		if err != nil {
			return err
		}
	} else {
		for i, le := range le {
			b := this.buckets[i]
			_, err = fmt.Fprintf(w, "%s_bucket{%s,le=\"%v\"} %d\n", name, labels, le, b)
			if err != nil {
				return err
			}
		}
		// +Inf is reported, even tho it's dup of "_sum"...
		_, err = fmt.Fprintf(w, "%s_bucket{%s,le=\"+Inf\"} %d\n", name, labels, this.count)
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(w, "%s_sum{%s} %f\n", name, labels, this.sum)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%s_count{%s} %d\n", name, labels, this.count)
		if err != nil {
			return err
		}
	}
	return nil
}
