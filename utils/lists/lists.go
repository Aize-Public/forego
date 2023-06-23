package lists

func Contains[T comparable](a []T, i T) bool {
	for _, e := range a {
		if e == i {
			return true
		}
	}
	return false
}

func AddUnique[T comparable](a []T, i T) []T {
	if Contains(a, i) {
		return a
	}
	return append(a, i)
}

func Copy[T any](a []T) (out []T) {
	out = make([]T, len(a))
	copy(out, a)
	return
}
