package prom

import (
	"io"
	"runtime"
	"runtime/debug"
	"time"
)

var Registry []Metric

type Metric interface {
	Print(io.Writer) error
}

func Register[T Metric](m T) T {
	Registry = append(Registry, m)
	return m
}

func init() {
	Register(&Gauge[int]{
		Name: "go_goroutines",
	}).SetFunc(func() int {
		return runtime.NumGoroutine()
	})

	Register(&Gauge[int]{
		Name:   "go_info",
		Labels: []string{"version"},
	}).SetFunc(func() int {
		return 1
	}, runtime.Version())

	Register(&Gauge[float64]{
		Name: "go_memstats_alloc_bytes",
	}).SetFunc(func() float64 {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		return float64(m.Alloc)
	})

	Register(&Custom{
		Name: "go_gc_duration_seconds",
		Type: "summary",
		Func: func() map[string]any {
			var s debug.GCStats
			s.PauseQuantiles = make([]time.Duration, 5)
			debug.ReadGCStats(&s)
			m := map[string]any{
				"go_gc_duration_seconds_count":            s.NumGC,
				"go_gc_duration_seconds_sum":              s.PauseTotal.Seconds(),
				`go_gc_duration_seconds{quantile="0"}`:    s.PauseQuantiles[0].Seconds(),
				`go_gc_duration_seconds{quantile="0.25"}`: s.PauseQuantiles[1].Seconds(),
				`go_gc_duration_seconds{quantile="0.5"}`:  s.PauseQuantiles[2].Seconds(),
				`go_gc_duration_seconds{quantile="0.75"}`: s.PauseQuantiles[3].Seconds(),
				`go_gc_duration_seconds{quantile="1"}`:    s.PauseQuantiles[4].Seconds(),
			}
			return m
		},
	})
}
