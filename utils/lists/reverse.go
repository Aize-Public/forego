package lists

// Shallow copy
func Copy[T any](a []T) []T {
	out := make([]T, len(a))
	for i := 0; i < len(a); i++ {
		out[i] = a[i]
	}
	return out
}

func Reverse[T any](a []T) {
	i := 0
	j := len(a) - 1
	for i < len(a)/2 {
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}
}
