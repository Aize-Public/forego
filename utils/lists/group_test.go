package lists

import (
	"github.com/Aize-Public/forego/test"
	"testing"
)

func TestSplit(t *testing.T) {

	{
		in := []int{}
		out := Copy(in)
		test.EqualsJSON(t, [][]int{{}}, Split(out, 1))
	}

	{
		in := []int{1, 2, 4, 8, 16}
		out := Copy(in)
		test.EqualsJSON(t, [][]int{{1, 2, 4, 8, 16}}, Split(out, 5))
		test.EqualsJSON(t, [][]int{{1, 2, 4, 8, 16}}, Split(out, 6))
		test.EqualsJSON(t, [][]int{{1, 2}, {4, 8}, {16}}, Split(out, 2))
		test.EqualsJSON(t, [][]int{{1, 2, 4}, {8, 16}}, Split(out, 4))
	}
	{
		in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
		out := Copy(in)
		test.EqualsJSON(t, [][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}}, Split(out, 13))
		test.EqualsJSON(t, [][]int{{1, 2, 3, 4, 5, 6, 7}, {8, 9, 10, 11, 12, 13}}, Split(out, 9))
	}

	{
		in := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
		out := Copy(in)
		test.EqualsJSON(t, [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12}}, Split(out, 6))
	}
}
