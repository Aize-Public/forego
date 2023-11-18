package lists_test

import (
	"testing"

	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/lists"
)

func TestSort(t *testing.T) {
	list := []uint{3, 7, 4, 0}
	t.Logf("%#v", list)
	lists.Sort(list, func(l, r uint) bool {
		return l > r
	})
	t.Logf("%#v", list)
	test.EqualsGo(t, []uint{7, 4, 3, 0}, list)
}

func TestSortFunc(t *testing.T) {
	list := []uint{3, 7, 4, 0}
	t.Logf("%#v", list)
	lists.SortFunc(list, func(v uint) int {
		return -int(v)
	})
	t.Logf("%#v", list)
	test.EqualsGo(t, []uint{7, 4, 3, 0}, list)
}
