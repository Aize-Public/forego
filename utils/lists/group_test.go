package lists

import (
	"testing"

	"github.com/Aize-Public/forego/test"
)

func TestSplit(t *testing.T) {
	c := test.Context(t)
	{
		in := []int{}
		out := Copy(in)
		test.EqualsJSON(c, [][]int{{}}, Split(out, 1))
	}

	{
		in := []int{1, 2, 4, 8, 16}
		out := Copy(in)
		test.EqualsJSON(c, [][]int{{1, 2, 4, 8, 16}}, Split(out, 5))
		test.EqualsJSON(c, [][]int{{1, 2, 4, 8, 16}}, Split(out, 6))
		test.EqualsJSON(c, [][]int{{1, 2}, {4, 8}, {16}}, Split(out, 2))
		test.EqualsJSON(c, [][]int{{1, 2, 4}, {8, 16}}, Split(out, 4))
	}
	{
		in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
		out := Copy(in)
		test.EqualsJSON(c, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}}, Split(out, 13))
		test.EqualsJSON(c, [][]int{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10, 11, 12, 13}}, Split(out, 9))
	}

	{
		in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		out := Copy(in)
		test.EqualsJSON(c, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12}}, Split(out, 6))
	}
}

func TestFlatten(t *testing.T) {
	c := test.Context(t)
	out := Flatten([][]int{{1}, {2, 3}, {}})
	test.EqualsJSON(c, `[1,2,3]`, out)
}
