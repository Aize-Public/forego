package lists

func Reverse[T any](a []T) {
	i, j := 0, len(a)-1
	for i < len(a)/2 {
		a[i], a[j] = a[j], a[i]
		i++
		j--
	}
}
