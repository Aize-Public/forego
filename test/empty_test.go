package test

import "testing"

func TestEmpty(t *testing.T) {
	empty(nil).true(t)

	empty("").true(t)
	empty("foo").false(t)

	empty([]int{}).true(t)
	empty(([]int)(nil)).true(t)
	empty([]int{1}).false(t)
	empty([]any{nil}).false(t)

	empty(map[int]int{}).true(t)
	empty(map[int]int{2: 3}).false(t)
}
