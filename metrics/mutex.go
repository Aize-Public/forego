package metrics

import (
	"fmt"
	"runtime"
	gosync "sync"
	"time"

	"github.com/Aize-Public/forego/metrics/prom"
	"github.com/Aize-Public/forego/sync"
)

var lock = prom.Register(&prom.Histogram{
	Name: "lock",
	Buckets: []float64{
		.0001, .00025, .0005,
		.001, .0025, .005,
		.01, .025, .05,
		.1, .25, .5,
		1, 2.5, 5},
	Labels: []string{"op", "src"},
})

type Mutex struct {
	m      gosync.Mutex
	t1     time.Time
	holder string
}

func (this *Mutex) Locker() *gosync.Mutex {
	return &this.m
}

func (this *Mutex) TryLock() bool {
	return this.m.TryLock()
}

func (this *Mutex) Lock() func() {
	// precompute src outside the lock
	_, file, line, _ := runtime.Caller(1)
	src := fmt.Sprintf("%s:%d", file, line)

	t0 := time.Now()
	if !sync.TryLock(10*time.Second, &this.m) {
		panic(fmt.Sprintf("lock timeout: %p (%+v)", &this.m, time.Since(t0)))
	}
	this.t1 = time.Now()
	lock.Observe(this.t1.Sub(t0).Seconds(), "acquire", this.holder)

	this.holder = src
	return this.Unlock
}

func (this *Mutex) Unlock() {
	h, t1 := this.holder, this.t1 // copy so they will be safe after unlock
	this.m.Unlock()
	lock.Observe(time.Since(t1).Seconds(), "release", h)
}
