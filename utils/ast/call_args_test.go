package ast_test

import (
	"errors"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/ast"
)

func testCall(c ctx.C, b any) (*ast.Call, string, error) {
	return ast.Caller(0)
}

func testAssign(c ctx.C, b any) string {
	return ast.Assignment(0, 1)
}

func TestArg(t *testing.T) {
	c := test.Context(t)

	i := 2
	call, _, err := testCall(c, i == 4/i)
	test.EqualsJSON(c, nil, err)
	test.EqualsJSON(c, "c", call.Args[0].Src)
	test.EqualsJSON(c, "i == 4/i", call.Args[1].Src)

	{
		err := errors.New("my error")
		src := testAssign(c, err)
		test.EqualsJSON(c, `errors.New("my error")`, src)
	}
}
