package ws

import (
	"github.com/Aize-Public/forego/utils/prom"
)

var Metrics = struct {
	Life    *prom.Histogram
	Request *prom.Histogram
	Gauge   *prom.Gauge
}{
	Life: prom.Register("ws_life", &prom.Histogram{
		Buckets: []float64{
			.1, .25, .5, 1, 2.5, 5, 10, 30,
			60, 60 * 5, 60 * 20,
			3600, 3600 * 3, 3600 * 8,
			86400, 86400 * 2, 86400 * 7,
		},
		Labels: []string{"path"},
	}),
	Request: prom.Register("ws_request", &prom.Histogram{
		Buckets: []float64{
			.001, .0025, .005,
			.01, .025, .05,
			.1, .25, .5,
			1, 2.5, 5,
			10, 20, 60,
		},
		Labels: []string{"path", "state"},
	}),
	Gauge: prom.Register("ws_gauge", &prom.Gauge{
		Labels: []string{"path"},
	}),
}
