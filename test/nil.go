package test

import (
	"fmt"
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

func Nil(t *testing.T, obj any) {
	t.Helper()
	isNil(obj).assignment(0, 1).true(t)
}

func NotNil(t *testing.T, obj any) {
	t.Helper()
	isNil(obj).assignment(0, 1).false(t)
}

func NoError(t *testing.T, err error) {
	t.Helper()
	if isNil(err).succeed {
		t.Logf("OK no error %s", ast.Assignment(0, 1))
	} else {
		t.Fatalf("FAIL %v", err)
	}
}

func Error(t *testing.T, err error) {
	t.Helper()
	isNil(err).assignment(0, 1).false(t)
}

func isNil(a any) res {
	switch a := a.(type) {
	case nil:
		return res{true, "nil"}
	default:
		return res{false, fmt.Sprint(a)}
	}
}
