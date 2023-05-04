package lists_test

import (
	"testing"

	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/lists"
)

func TestReverse(t *testing.T) {
	{
		in := []int{1, 2, 4, 8, 16}
		out := lists.Copy(in)
		lists.Reverse(out)
		test.EqualsJSON(t, []int{1, 2, 4, 8, 16}, in)
		test.EqualsJSON(t, []int{16, 8, 4, 2, 1}, out)
	}
}
