package sync

import (
	"sync"
	"time"
)

// this to avoid importing 2 sync packages
type Mutex struct {
	sync.Mutex
}

// attempt to obtain a lock to the given Mutex, or return false if timeout
func TryLock(d time.Duration, m sync.Locker) bool {
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
