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
	m      sync.Mutex
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
