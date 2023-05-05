package test

import (
	"github.com/Aize-Public/forego/utils/ast"
)

func Assert(t T, cond bool) {
	t.Helper()
	call, err := ast.Caller(0)
	if err != nil {
		t.Logf("can't parse ast: %v", err)
		if cond {
			t.Logf("ok")
		} else {
			t.Fatalf("FAIL")
		}
		return
	}
	if cond {
		t.Logf("ok: %s", call.Args[1])
	} else {
		t.Fatalf("FAIL: %s", call.Args[1])
	}
}
