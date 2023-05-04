package test

import (
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

func Assert(t *testing.T, cond bool) {
	t.Helper()
	call, err := ast.Caller(0)
	if err != nil {
		t.Logf("can't parse ast: %v", err)
		if cond {
			t.Logf("ok")
		} else {
			t.Fatal("FAIL")
		}
		return
	}
	if cond {
		t.Logf("ok: %s", call.Args[1])
	} else {
		t.Fatalf("FAIL: %s", call.Args[1])
	}
}
