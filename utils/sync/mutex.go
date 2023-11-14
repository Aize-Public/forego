package sync

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/Aize-Public/forego/ctx/log"
)

type Mutex struct {
	sync.Mutex
}

type MetricsMutex struct {
	m      Mutex
	lock   time.Time
	holder string
}

func (m *MetricsMutex) Lock() {
	_, file, line, _ := runtime.Caller(1)
	holder := fmt.Sprintf("%s:%d", file, line)
	t0 := time.Now()
	m.m.Lock()
	m.lock = time.Now()
	m.holder = holder
	// TODO metrics instead of logs
	log.Debugf(nil, "lock() in %v at %s", m.lock.Sub(t0), holder)
}

func (m *MetricsMutex) Unlock() {
	dt := time.Since(m.lock)
	holder := m.holder
	m.m.Unlock()
	// TODO metrics instead of logs
	log.Debugf(nil, "unlock() in %v at %s", dt, holder)
}

// attempt to obtain a lock to the given Mutex, or return false if timeout
func (m *Mutex) TryLock(d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	// TODO maybe we can use a channel internally instead of a chan + lock?
	ch := make(chan struct{})
	defer close(ch)
	go func() {
		m.Lock()
		_, ok := <-ch
		if !ok {
			m.Unlock()
		}
	}()
	select {
	case ch <- struct{}{}:
		return true
	case <-t.C:
		return false
	}
}
