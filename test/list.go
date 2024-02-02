package test

import (
	"fmt"
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

// shift the first element of the given array, or Fail()
func Shift[T any](t *testing.T, list *[]T) T {
	t.Helper()
	el, res := shift(list)
	res.prefix("Unshift(%s)", ast.Assignment(0, 1)).true(t)
	return el
}

// pop the last element of the given array, or Fail()
func Pop[T any](t *testing.T, list *[]T) T {
	t.Helper()
	el, res := pop(list)
	res.prefix("Pop(%s)", ast.Assignment(0, 1)).true(t)
	return el
}

func shift[T any](list *[]T) (t T, r res) {
	switch len(*list) {
	case 0:
		return t, res{false, "empty"}
	default:
		head := (*list)[0]
		*list = (*list)[1:]
		return head, res{true, fmt.Sprint(head)}
	}
}

func pop[T any](list *[]T) (t T, r res) {
	switch len(*list) {
	case 0:
		return t, res{false, "empty"}
	default:
		tail := (*list)[len(*list)-1]
		*list = (*list)[:len(*list)-1]
		return tail, res{true, fmt.Sprint(tail)}
	}
}
