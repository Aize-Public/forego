package utils_test

import (
	"testing"

	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils"
)

func TestMap(t *testing.T) {
	m := utils.NewOrderedMap[string, int]()
	test.Assert(t, len(m.Keys()) == 0)
	test.Assert(t, m.Get("one").Found == false)

	t.Logf("adding one")
	test.Assert(t, m.Put("one", 1).Found == false)
	test.Assert(t, m.Get("one").Found == true)
	test.Assert(t, m.Get("one").Value == 1)
	test.Assert(t, len(m.Keys()) == 1)

	t.Logf("adding two")
	test.Assert(t, m.Put("two", 2).Found == false)
	test.Assert(t, m.Get("two").Found == true)
	test.Assert(t, m.Get("two").Value == 2)
	test.EqualsJSON(t, []string{"one", "two"}, m.Keys()) // orders preserved

	t.Logf("changing one won't change order")
	test.Assert(t, m.Put("one", -1).Found == true)
	test.Assert(t, m.Get("one").Found == true)
	test.Assert(t, m.Get("one").Value == -1)
	test.EqualsJSON(t, []int{-1, 2}, m.Values()) // orders preserved

	t.Logf("adding more...")
	test.Assert(t, m.Put("life", 42).Found == false)
	test.Assert(t, m.Put("everything", 42).Found == false)
	test.EqualsGo(t, 4, m.Len())

	test.Assert(t, m.Delete("one").Found == true)
	test.EqualsJSON(t, []string{"two", "life", "everything"}, m.Keys()) // orders preserved

	test.Assert(t, m.Put("one", 1).Found == false)
	test.Assert(t, m.Get("one").Found == true)
	test.Assert(t, m.Get("one").Value == 1)
	test.EqualsJSON(t, []string{"two", "life", "everything", "one"}, m.Keys()) // added last after delete
}
