package http

import (
	"fmt"
	"time"

	"github.com/Aize-Public/forego/utils/prom"
)

var Metrics = struct {
	Request *prom.Histogram
}{
	Request: prom.Register("http_request", &prom.Histogram{
		Buckets: prom.DefaultBuckets,
		Labels:  []string{"method", "path", "code"},
	}),
	// TODO add more, like active requests gauges, open connections, websockets...
}

type metric struct {
	Method string
	Path   string
	Code   int
}

func (m metric) observe(d time.Duration) {
	Metrics.Request.Observe(d.Seconds(), m.Method, m.Path, fmt.Sprint(m.Code))
}
