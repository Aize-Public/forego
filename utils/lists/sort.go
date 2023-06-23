package lists

import (
	"fmt"
	"sort"
)

func Sort[T any](a []T, less func(T, T) bool) {
	sort.Slice(a, func(l, r int) bool {
		return less(a[l], a[r])
	})
}

func SortWithSprint[T any](a []T) {
	sort.Slice(a, func(l, r int) bool {
		return fmt.Sprint(a[l]) < fmt.Sprint(a[r])
	})
}

type comparable interface {
	int | float64 | string
}

func SortFunc[T any, S comparable](a []T, score func(T) S) {
	sort.Slice(a, func(l, r int) bool {
		al := score(a[l])
		ar := score(a[r])
		return al < ar
	})
}
