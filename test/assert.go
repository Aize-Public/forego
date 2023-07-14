package test

import (
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

// var ok = "✅  "
var ok = "  ✔ "
var ko = "❌  "

func OK(t *testing.T, f string, args ...any) {
	t.Helper()
	t.Logf(ok+f, args...)
}

func Fail(t *testing.T, f string, args ...any) {
	t.Helper()
	t.Fatalf(ko+f, args...)
}

func Assert(t *testing.T, cond bool) {
	t.Helper()
	if cond {
		OK(t, "%s", stringy{ast.Assignment(0, 1)})
	} else {
		Fail(t, "%s", stringy{ast.Assignment(0, 1)})
	}
}
