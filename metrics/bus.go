package metrics

import (
	"time"

	"github.com/Aize-Public/forego/metrics/prom"
)

type BusPublish struct {
	Topic string
}

var busPublish = prom.Register(&prom.Counter{
	Name:   "bus_publish",
	Labels: []string{"topic"},
})

func (m BusPublish) Observe() {
	busPublish.Observe(1, m.Topic)
}

type BusLag struct {
	Label string
}

var busLag = prom.Register(&prom.Histogram{
	Name:    "bus_lag",
	Buckets: []float64{.001, .0025, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 30, 60, 60 * 5}, // 1ms .. 5'
	Labels:  []string{"label"},
})

func (m BusLag) Observe(lag time.Duration) {
	busLag.Observe(lag.Seconds(), m.Label)
}
