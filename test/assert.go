package test

import (
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

func Assert(t *testing.T, cond bool) {
	t.Helper()
	if cond {
		t.Logf("ok %s", ast.Assignment(0, 1))
	} else {
		t.Fatalf("fail %s", ast.Assignment(0, 1))
	}
}
