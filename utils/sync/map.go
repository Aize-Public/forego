package sync

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Map[K any, T any] struct {
	m sync.Map
}

func (this *Map[K, T]) MarshalJSON() (out []byte, err error) {
	m := map[string]json.RawMessage{}
	this.Range(func(k K, v T) bool {
		ks := fmt.Sprintf("%v", k)
		m[ks], err = json.Marshal(v)
		return err == nil
	})
	if err != nil {
		return
	}
	return json.Marshal(m)
}

func (this *Map[K, T]) Get(key K) T {
	val, _ := this.m.Load(key)
	if val == nil {
		var zero T
		return zero
	}
	return val.(T)
}

func (this *Map[K, T]) IsEmpty() bool {
	empty := true
	this.m.Range(func(_, _ any) bool {
		empty = false
		return false
	})
	return empty
}

func (this *Map[K, T]) Load(key K) (T, bool) {
	val, ok := this.m.Load(key)
	if val == nil {
		var zero T
		return zero, ok
	}
	return val.(T), ok
}

func (this *Map[K, T]) GetOrStore(key K, val T) T {
	out, _ := this.m.LoadOrStore(key, val)
	if out == nil {
		var zero T
		return zero
	}
	return out.(T)
}

func (this *Map[K, T]) LoadOrStore(key K, val T) (T, bool) {
	out, loaded := this.m.LoadOrStore(key, val)
	if out == nil {
		var zero T
		return zero, loaded
	}
	return out.(T), loaded
}

func (this *Map[K, T]) Store(key K, val T) {
	this.m.Store(key, val)
}

func (this *Map[K, T]) Range(f func(key K, val T) bool) {
	this.m.Range(func(key, val any) bool {
		return f(key.(K), val.(T))
	})
}

func (this *Map[K, T]) Delete(key K) {
	this.m.Delete(key)
}

func (this *Map[K, T]) LoadAndDelete(key K) (T, bool) {
	v, loaded := this.m.LoadAndDelete(key)
	if v == nil {
		var zero T
		return zero, loaded
	}
	return v.(T), loaded
}
