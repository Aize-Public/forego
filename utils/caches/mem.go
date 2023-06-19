package caches

import (
	"sync"

	"github.com/Aize-Public/forego/ctx"
)

type memCache[K comparable, V any] struct {
	miss func(ctx.C, K) (V, error)
	lock sync.Mutex
	m    map[K]V
}

func Mem[K comparable, V any](c ctx.C, miss func(ctx.C, K) (V, error)) Cache[K, V] {
	return &memCache[K, V]{
		miss: miss,
		m:    map[K]V{},
	}
}

func (this *memCache[K, V]) Invalidate(c ctx.C, k K) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.m, k)
}

func (this *memCache[K, V]) Get(c ctx.C, k K) (V, HitMiss, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	v, exists := this.m[k]
	if exists {
		return v, HIT, nil
	}
	v, err := this.miss(c, k)
	if err != nil {
		return v, ERR, err
	}
	this.m[k] = v
	return v, MISS, nil
}
