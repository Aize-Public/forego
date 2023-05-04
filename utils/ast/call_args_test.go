package ast_test

import (
	"errors"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/ast"
)

func testCall(c ctx.C, b any) (*ast.Call, error) {
	return ast.Caller(0)
}

func TestArg(t *testing.T) {
	c := ctx.TODO()
	i := 2
	call, err := testCall(c, i == 4/i)
	test.EqualsJSON(t, nil, err)
	test.EqualsJSON(t, "c", call.Args[0].Src)
	test.EqualsJSON(t, "i == 4/i", call.Args[1].Src)

	{
		err := errors.New("my error")
		call, err_ := testCall(c, err)
		test.EqualsJSON(t, nil, err_)

		test.EqualsJSON(t, "c", call.Args[0].Src)

		test.EqualsJSON(t, "err", call.Args[1].Src)
		test.EqualsJSON(t, `errors.New("my error")`, call.Args[1].Assignment)
	}
}
