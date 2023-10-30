package lists

func Split[T any](in []T, maxSplitSize int) [][]T {
	if maxSplitSize <= 0 {
		panic("maxSplitSize must be greater than 0")
	}
	if len(in) <= maxSplitSize {
		return [][]T{in}
	}

	buckets := (len(in) + maxSplitSize - 1) / maxSplitSize
	out := make([][]T, buckets)
	chunkSize := (len(in) + buckets - 1) / buckets

	for i := 0; i < buckets; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(in) {
			end = len(in)
		}
		out[i] = in[start:end]
	}

	return out
}

func Flatten[T any](in [][]T) []T {
	var out []T
	for _, group := range in {
		out = append(out, group...)
	}
	return out
}
