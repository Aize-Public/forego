package maps

type PairList[K comparable, V any] []Pair[K, V]

func (this PairList[K, V]) Keys() []K {
	out := make([]K, len(this))
	for i, p := range this {
		out[i] = p.Key
	}
	return out
}

func (this PairList[K, V]) Values() []V {
	out := make([]V, len(this))
	for i, p := range this {
		out[i] = p.Value
	}
	return out
}

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func Pairs[K comparable, V any](in map[K]V) (out PairList[K, V]) {
	for k, v := range in {
		out = append(out, Pair[K, V]{
			Key:   k,
			Value: v,
		})
	}
	return out
}
