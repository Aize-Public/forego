package lists

// Shallow copy
func Copy[T any](a []T) (out []T) {
	out = make([]T, len(a))
	copy(out, a)
	return
}

func Reverse[T any](a []T) {
	i, j := 0, len(a)-1
	for i < len(a)/2 {
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}
}

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
