package caches

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/Aize-Public/forego/ctx"
)

type LRU[K comparable, V any] struct {
	l    sync.Mutex
	heap *lruHeap[K, V]

	MaxSize int
	curSize int
	// called when there is a cache miss, allowing the implementation to make room for the new entry
	OnMiss func(c ctx.C, k K) (V, int, error)
}

func NewLRU[K comparable, V any](c ctx.C, maxSize int,
	onMiss func(ctx.C, K) (v V, elementSize int, err error),
) Cache[K, V] {
	return &LRU[K, V]{
		heap: &lruHeap[K, V]{
			index: map[K]int{},
		},
		OnMiss:  onMiss,
		MaxSize: maxSize,
	}
}

func (this *LRU[K, V]) Keys() []K {
	this.l.Lock()
	defer this.l.Unlock()

	keys := make([]K, 0, len(this.heap.list))
	for _, e := range this.heap.list {
		keys = append(keys, e.key)
	}
	return keys
}

func (this *LRU[K, V]) Invalidate(c ctx.C, k K) {
	this.l.Lock()
	defer this.l.Unlock()

	i, exists := this.heap.index[k]
	if !exists {
		return
	}
	this.heap.list[i].last = time.Time{} // make it the oldest
	heap.Fix(this.heap, i)               // make it go to the front

	entry := heap.Pop(this.heap).(lruEntry[K, V]) // remove from the heap (NOTE(oha): can't see how the element would not be in the front of the heap)
	delete(this.heap.index, entry.key)            // remove from the index
	this.curSize -= entry.size
}

func (this *LRU[K, V]) Get(c ctx.C, k K) (V, HitMiss, error) {
	this.l.Lock()
	defer this.l.Unlock()

	i, exists := this.heap.index[k]
	if exists {
		entry := this.heap.list[i]
		entry.last = time.Now()
		heap.Fix(this.heap, i)
		return entry.val, HIT, nil
	}

	// TODO(oha): this is wrong, we shouldn't hold a lock onMiss
	// instead we should use some more magic to make only the same key block
	v, size, err := this.OnMiss(c, k)
	if err != nil {
		return v, ERR, err
	}
	if size == 0 { // no cache
		return v, NO, nil
	}

	entry := lruEntry[K, V]{
		key:  k,
		val:  v,
		last: time.Now(),
		size: size,
	}
	heap.Push(this.heap, entry)
	this.curSize += entry.size

	for this.curSize > this.MaxSize {
		entry := heap.Pop(this.heap).(lruEntry[K, V])
		this.curSize -= entry.size
		//log.Debugf(c, "evicted %+v", entry.key)
	}
	return v, MISS, nil
}

func (this *LRU[K, V]) String() string {
	if len(this.heap.list) == 0 {
		return fmt.Sprintf("LRU{%d/%d, empty}", this.curSize, this.MaxSize)
	}
	return fmt.Sprintf("LRU{%d/%d, %d entries, oldest %v}",
		this.curSize, this.MaxSize,
		len(this.heap.list), this.heap.list[0].key)
}

// implements heap.Interface{}

type lruHeap[K comparable, V any] struct {
	list  []lruEntry[K, V]
	index map[K]int
}

type lruEntry[K comparable, V any] struct {
	key  K
	val  V
	last time.Time
	size int
	m    sync.Mutex
}

func (this *lruHeap[K, V]) Len() int {
	return len(this.list)
}

func (this *lruHeap[K, V]) Less(i, j int) bool {
	return this.list[i].last.Before(this.list[j].last)
}

func (this *lruHeap[K, V]) Swap(i, j int) {
	// swap
	this.list[i], this.list[j] = this.list[j], this.list[i]
	// update the index
	this.index[this.list[i].key] = i
	this.index[this.list[j].key] = j
}

func (this *lruHeap[K, V]) Pop() any {
	entry := this.list[len(this.list)-1]
	this.list = this.list[0 : len(this.list)-1]
	delete(this.index, entry.key)
	return entry
}

func (this *lruHeap[K, V]) Push(v any) {
	entry := v.(lruEntry[K, V])
	this.index[entry.key] = len(this.list) // update the index with the current position
	this.list = append(this.list, entry)
}
