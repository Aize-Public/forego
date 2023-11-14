package prom

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/Aize-Public/forego/utils/sync"
)

var registry = struct {
	sync.Mutex
	metrics map[string]Metric
}{
	metrics: map[string]Metric{},
}

func Register[T Metric](name string, m T) T {
	registry.Lock()
	defer registry.Unlock()
	old := registry.metrics[name]
	if old != nil {
		panic(fmt.Sprintf("duplicate metrics %q", name))
	}
	registry.metrics[name] = m
	return m
}

func Range(f func(name string, m Metric) error) error {
	out := map[string]Metric{}
	registry.Lock()
	for name, m := range registry.metrics {
		out[name] = m
	}
	registry.Unlock()
	for name, m := range out {
		err := f(name, m)
		if err != nil {
			return err
		}
	}
	return nil
}

type Metric interface {
	Print(string, io.Writer) error
}

func init() {
	Register("go_routines", &Gauge{}).SetFunc(func() float64 {
		return float64(runtime.NumGoroutine())
	})

	Register("go_info", &Gauge{
		Labels: []string{"version"},
	}).SetFunc(func() float64 {
		return 1.0
	}, runtime.Version())

	Register("go_memstats_alloc_bytes", &Gauge{}).SetFunc(func() float64 {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		return float64(m.Alloc)
	})

	Register("go_gc_duration_seconds", &Custom{
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
