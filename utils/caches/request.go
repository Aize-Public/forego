package caches

import (
	"sync"
	"time"

	"github.com/Aize-Public/forego/ctx"
)

type RequestCache[K comparable, V any] struct {
	cache  Cache[K, *entry[V]]
	onMiss func(c ctx.C, key K) (V, time.Duration, error)
}

type entry[V any] struct {
	m      sync.Mutex
	V      V
	expire time.Time
	err    error
}

func NewRequestLRU[K comparable, V any](c ctx.C, maxSize int, onMiss func(c ctx.C, k K) (V, time.Duration, error)) *RequestCache[K, V] {
	return &RequestCache[K, V]{
		onMiss: onMiss,
		cache: NewLRU(c, maxSize, func(c ctx.C, k K) (*entry[V], int, error) {
			return &entry[V]{}, 1, nil
		}),
	}
}

func (this *RequestCache[K, V]) Invalidate(c ctx.C, key K) {
	this.cache.Invalidate(c, key)
}

func (this *RequestCache[K, V]) Request(c ctx.C, key K) (V, HitMiss, error) {
	var zero V
	e, hit, err := this.cache.Get(c, key)
	if err != nil {
		return zero, hit, err
	}

	// only one should have access to this entry
	e.m.Lock()
	defer e.m.Unlock()

	//log.Debugf(c, "ReqCache[%v] %v expired %v ago", key, hit, time.Since(e.expire))
	// not a miss, and not expired
	if hit != MISS && time.Until(e.expire) >= 0 {
		if e.err != nil {
			return e.V, HIT, e.err
		} else {
			return e.V, HIT, nil
		}
	}

	// FETCH
	v, ttl, err := this.onMiss(c, key)
	if ttl > 0 {
		e.err = err
		e.V = v
		e.expire = time.Now().Add(ttl)
	} else {
		e.expire = time.Time{} // zero so it looks expired
	}

	if err != nil {
		return v, ERR, err
	} else {
		return v, MISS, nil
	}
}
