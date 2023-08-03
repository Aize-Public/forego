package test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

func IsType(t *testing.T, expect any, got any) {
	t.Helper()
	isType(expect, got).prefix("%s == %s", Quote(ast.Assignment(0, 1)), Quote(ast.Assignment(0, 2))).true(t)
}

func isType(e, g any) res {
	et := reflect.TypeOf(e)
	gt := reflect.TypeOf(e)
	if et == gt {
		return res{true, et.String()}
	} else {
		return res{false, fmt.Sprintf("%T != %T", e, g)}
	}
}
