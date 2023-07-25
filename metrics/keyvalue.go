package metrics

import (
	"time"

	"github.com/Aize-Public/forego/metrics/prom"
)

type KeyValue struct {
	Op  string // "tx" for the whole transaction, "exec" for exec statements, "query" or "queryAll" (note queryAll are probably pointless)
	Src string
}

var keyvalue = prom.Register(&prom.Histogram{
	Name:    "keyvalue",
	Buckets: []float64{.001, .0025, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30, 60, 60 * 5}, // 1ms .. 5'
	Labels:  []string{"op", "src"},
})

func (m KeyValue) ObserveSince(start time.Time) {
	keyvalue.Observe(time.Since(start).Seconds(), m.Op, m.Src)
}
