package lists

import "sync"

type Set[T comparable] struct {
	m sync.Map
}

func (this *Set[T]) Add(t T) bool {
	_, exists := this.m.LoadOrStore(t, true)
	return !exists
}

func (this *Set[T]) Has(t T) bool {
	_, exists := this.m.Load(t)
	return exists
}

func (this *Set[T]) Remove(t T) bool {
	return this.m.CompareAndDelete(t, true)
}

func (this *Set[T]) All() (list []T) {
	this.m.Range(func(k, _ any) bool {
		list = append(list, k.(T))
		return true
	})
	return
}
