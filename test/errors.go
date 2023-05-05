package test

import "github.com/Aize-Public/forego/utils/ast"

func NoError(t T, err error) {
	t.Helper()

	if err == nil {
		t.Logf("ok: %s", ast.Assignment(1, 1))
	} else {
		t.Fatalf("ERROR: %s: %v", ast.Assignment(1, 1), err)
	}
}
