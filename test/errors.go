package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/utils/ast"
)

func NoError(t *testing.T, err error) {
	t.Helper()
	if isNil(err).succeed {
		OK(t, "no error: %s", stringy{ast.Assignment(0, 1)})
	} else {
		var cErr ctx.Error
		if errors.As(err, &cErr) {
			Fail(t, "%v\n\t%s", err, strings.Join(cErr.Stack, "\n\t"))
		} else {
			Fail(t, "%v", err)
		}
	}
}

func Error(t *testing.T, err error) {
	t.Helper()
	isNil(err).assignment(0, 1).false(t)

	if isNil(err).succeed {
		Fail(t, "expected error: ", stringy{ast.Assignment(0, 1)})
	} else {
		var cErr ctx.Error
		if errors.As(err, &cErr) {
			OK(t, "%v\n\t%s", err, strings.Join(cErr.Stack, "\n\t"))
		} else {
			OK(t, "%v", err)
		}
	}
}
