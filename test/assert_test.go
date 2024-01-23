package test

import (
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

func testAssignment(t *testing.T, cond bool) {
	c := Context(t)
	a := ast.Assignment(0, 1)
	t.Logf("assignment: %s", a)
	ContainsJSON(c, a, "42")
}
func TestAssert(t *testing.T) {
	yes := 42 > 7
	Assert(t, yes)
	testAssignment(t, yes)
}
