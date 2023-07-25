package sync

import (
	"errors"
	"sync"
)

// make the non-generic sync.Map into a sync.Map[K,V]
type Map[K comparable, V any] struct {
	m sync.Map
}

func (this *Map[K, V]) Get(k K) (v V) {
	vv, ok := this.m.Load(k)
	if ok {
		return vv.(V)
	} else {
		return
	}
}

func (this *Map[K, V]) Load(k K) (v V, exists bool) {
	vv, ok := this.m.Load(k)
	if ok {
		return vv.(V), true
	} else {
		return
	}
}

func (this *Map[K, V]) GetOrStore(k K, v V) V {
	vv, _ := this.m.LoadOrStore(k, v)
	return vv.(V)
}

// used to interupt operations without really returning an error
var EOD = errors.New("EOD")

// range over all the pairs, until an error is returned
// use sync.EOD to quite the loop without giving an error
func (this *Map[K, V]) Range(f func(k K, v V) error) error {
	var err error
	this.m.Range(func(k, v any) bool {
		err = f(k.(K), v.(V))
		return err == nil
	})
	if err == EOD {
		return nil
	}
	return err
}

func (this *Map[K, V]) Delete(k K) {
	this.m.Delete(k)
}
func (this *Map[K, V]) Store(k K, v V) {
	this.m.Store(k, v)
}
