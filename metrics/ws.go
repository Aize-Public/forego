package metrics

import (
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/metrics/prom"
)

type WS struct {
	Path string
}

type WSRequest struct {
	Path    string
	Subpath string
}

// websocket life
var ws = prom.Register(&prom.Histogram{
	Name: "ws",
	Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10, 30,
		60, 60 * 5, 60 * 20,
		3600, 3600 * 3, 3600 * 8,
		86400, 86400 * 2, 86400 * 7,
	},
	Labels: []string{"path"},
})

// request within a websocket // NOTE(oha) we might want to just use api-metrics here instead
var wsRequest = prom.Register(&prom.Histogram{
	Name: "ws_request",
	Buckets: []float64{
		.001, .0025, .005,
		.01, .025, .05,
		.1, .25, .5,
		1, 2.5, 5,
		10, 20, 60,
	},
	Labels: []string{"path", "subpath", "exit"},
})
var wsGauge = prom.Register(&prom.Counter{
	Name:   "ws_gauge",
	Labels: []string{"path"},
	Gauge:  true,
})

func (m WS) Start() End {
	start := time.Now()
	wsGauge.Observe(1, m.Path)
	return func(c ctx.C) {
		wsGauge.Observe(-1, m.Path)
		dt := time.Since(start)
		ws.Observe(dt.Seconds(), m.Path)
		log.Infof(c, "WS %q End() after %v", m.Path, dt)
	}
}

type End func(ctx.C)

func (this End) End(c ctx.C) { this(c) }

func (m WSRequest) Since(since time.Time, err error) {
	if err == nil {
		wsRequest.Observe(time.Since(since).Seconds(), m.Path, m.Subpath, "ok")
	} else {
		wsRequest.Observe(time.Since(since).Seconds(), m.Path, m.Subpath, "err")
	}
}
