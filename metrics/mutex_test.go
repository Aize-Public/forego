package metrics_test

import (
	"os"
	"runtime/pprof"
	"testing"

	"github.com/Aize-Public/forego/metrics"
	"github.com/Aize-Public/forego/test"
)

func TestMutext(t *testing.T) {
	f, err := os.Create("/tmp/cpu.prof")
	if err != nil {
		panic(err)
	}
	_ = pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	for i := 0; i < 10000; i++ {
		var m metrics.Mutex
		m.Lock()
		test.Assert(t, false == m.TryLock())
		m.Unlock()

		test.Assert(t, true == m.TryLock())
		test.Assert(t, false == m.TryLock())
		m.Unlock()
	}
}
