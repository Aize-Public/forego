package utils

type OrderedMap[K comparable, V any] struct {
	list  []*Pair[K, V]
	byKey map[K]int
}

type Pair[K any, V any] struct {
	Key   K
	Value V
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{
		byKey: map[K]int{},
	}
}

type MapEntry[V any] struct {
	Value V
	Found bool
}

func (this OrderedMap[K, V]) Get(key K) MapEntry[V] {
	i, ok := this.byKey[key]
	if ok {
		return MapEntry[V]{this.list[i].Value, true}
	} else {
		return MapEntry[V]{}
	}
}

func (this *OrderedMap[K, V]) Put(key K, value V) MapEntry[V] {
	i, ok := this.byKey[key]
	if ok {
		old := this.list[i].Value
		this.list[i].Value = value
		return MapEntry[V]{old, true}
	} else {
		i := len(this.list)
		this.list = append(this.list, &Pair[K, V]{key, value})
		this.byKey[key] = i
		return MapEntry[V]{}
	}
}

// remove the pair and forget about the order
func (this *OrderedMap[K, V]) Delete(key K) MapEntry[V] {
	i, ok := this.byKey[key]
	if ok {
		old := this.list[i].Value
		delete(this.byKey, key)
		this.list[i] = nil
		return MapEntry[V]{old, true}
	} else {
		return MapEntry[V]{}
	}
}

func (this OrderedMap[K, V]) Range(f func(K, V) error) error {
	for _, pair := range this.list {
		if pair == nil {
			continue
		}
		err := f(pair.Key, pair.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this OrderedMap[K, V]) Len() int {
	len := 0
	this.Range(func(k K, v V) error {
		len++
		return nil
	})
	return len
}

func (this OrderedMap[K, V]) Keys() []K {
	list := []K{}
	this.Range(func(k K, v V) error {
		list = append(list, k)
		return nil
	})
	return list
}

func (this OrderedMap[K, V]) Values() []V {
	list := []V{}
	this.Range(func(k K, v V) error {
		list = append(list, v)
		return nil
	})
	return list
}
