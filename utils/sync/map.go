package sync

import (
	"sync"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/enc"
)

type Map[K comparable, T any] struct {
	m sync.Map
}

var _ enc.Marshaler = (*Map[string, any])(nil)
var _ enc.Unmarshaler = (*Map[string, any])(nil)

func (this *Map[K, V]) Map() map[K]V {
	m := map[K]V{}
	_ = this.RangeErr(func(k K, v V) error {
		m[k] = v
		return nil
	})
	return m
}

func (this *Map[K, V]) MarshalNode(c ctx.C) (n enc.Node, err error) {
	m := this.Map()
	return enc.Marshal(c, m)
}

func (this *Map[K, V]) UnmarshalNode(c ctx.C, n enc.Node) error {
	var m map[K]V
	err := enc.Unmarshal(c, n, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		this.Store(k, v)
	}
	return nil
}

func (this *Map[K, V]) Get(key K) V {
	val, _ := this.m.Load(key)
	if val == nil {
		var zero V
		return zero
	}
	return val.(V)
}

func (this *Map[K, V]) IsEmpty() bool {
	empty := true
	this.m.Range(func(_, _ any) bool {
		empty = false
		return false
	})
	return empty
}

func (this *Map[K, V]) Load(key K) (V, bool) {
	val, ok := this.m.Load(key)
	if val == nil {
		var zero V
		return zero, ok
	}
	return val.(V), ok
}

func (this *Map[K, V]) GetOrStore(key K, val V) V {
	out, _ := this.m.LoadOrStore(key, val)
	if out == nil {
		var zero V
		return zero
	}
	return out.(V)
}

func (this *Map[K, V]) LoadOrStore(key K, val V) (V, bool) {
	out, loaded := this.m.LoadOrStore(key, val)
	if out == nil {
		var zero V
		return zero, loaded
	}
	return out.(V), loaded
}

func (this *Map[K, V]) Store(key K, val V) {
	this.m.Store(key, val)
}

func (this *Map[K, V]) Range(f func(key K, val V) bool) {
	this.m.Range(func(key, val any) bool {
		return f(key.(K), val.(V))
	})
}

// like Range, but using errors not bools
func (this *Map[K, V]) RangeErr(f func(key K, val V) error) error {
	var err error
	this.m.Range(func(key, val any) bool {
		err = f(key.(K), val.(V))
		return err == nil
	})
	return err
}

func (this *Map[K, V]) Delete(key K) {
	this.m.Delete(key)
}

func (this *Map[K, V]) LoadAndDelete(key K) (V, bool) {
	v, loaded := this.m.LoadAndDelete(key)
	if v == nil {
		var zero V
		return zero, loaded
	}
	return v.(V), loaded
}
