package lists_test

import (
	"testing"

	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/lists"
)

func TestUnique(t *testing.T) {
	in := []int{}
	in = lists.AddUnique(in, 1)
	in = lists.AddUnique(in, 2)
	in = lists.AddUnique(in, 1)
	in = lists.AddUnique(in, 1)
	test.EqualsJSON(t, `[1,2]`, in)
}
