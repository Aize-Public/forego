package ctx_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/test"
)

func TestError(t *testing.T) {
	c := ctx.TODO()
	var err error
	err = io.EOF
	t.Logf("err: %T %v", err, err)

	err = ctx.NewErrorf(c, "wrap: %w", err)
	var x ctx.Error
	test.Assert(t, errors.As(err, &x))
	t.Logf("err: %T %v", err, err)

	stack := x.Stack[0]
	t.Logf("stack: %+v", stack)

	err = ctx.WrapError(c, err)
	t.Logf("err: %T %v", err, err)

	err = fmt.Errorf("wrap more: %w", err)
	t.Logf("err: %T %v", err, err)

	err = ctx.NewErrorf(c, "new: %w", err)
	t.Logf("err: %T %v", err, err)
	test.Error(t, err)

	var cerr ctx.Error
	ok := errors.As(err, &cerr)
	test.Assert(t, ok)
	test.Error(t, cerr)
	t.Logf("err: %s", err.Error())

	test.EqualsStr(t, stack, cerr.Stack[0])
}
