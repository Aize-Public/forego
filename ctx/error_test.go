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
	err = ctx.Errorf(c, "wrap: %w", err)
	err = ctx.Error(c, err)
	err = fmt.Errorf("wrap more: %w", err)
	test.Error(t, err)

	var cerr ctx.Err
	ok := errors.As(err, &cerr)
	test.Assert(t, ok)
	test.Error(t, cerr)
	t.Logf("stack: %+v", cerr.Stack)
}
