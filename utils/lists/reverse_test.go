package lists_test

import (
	"testing"

	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/lists"
)

func TestReverse(t *testing.T) {
	c := test.Context(t)
	{
		in := []int{1, 2, 4, 8, 16}
		out := lists.Copy(in)
		lists.Reverse(out)
		test.EqualsJSON(c, []int{1, 2, 4, 8, 16}, in)
		test.EqualsJSON(c, []int{16, 8, 4, 2, 1}, out)
	}
}
