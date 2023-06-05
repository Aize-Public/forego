package test

import (
	"fmt"
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

func EqualsStr(t *testing.T, expect, got string) {
	t.Helper()
	equal(expect, got).prefix("%s == %s", Quote(ast.Assignment(0, 1)), Quote(ast.Assignment(0, 2))).true(t)
}

func NotEqualsStr(t *testing.T, expect, got string) {
	t.Helper()
	equal(expect, got).prefix("%s == %s", Quote(ast.Assignment(0, 1)), Quote(ast.Assignment(0, 2))).false(t)
}

func equal(e, g string) res {
	if e == g {
		return res{true, e}
	} else {
		return res{false, fmt.Sprintf("%s != %s", Quote(e), Quote(g))}
	}
}

// compare using "%#v"
func EqualsGo(t *testing.T, expect, got any) {
	t.Helper()
	equalGo(expect, got).prefix("%s == %s", Quote(ast.Assignment(0, 1)), Quote(ast.Assignment(0, 2))).true(t)
}

// compare using "%#v"
func NotEqualsGo(t *testing.T, expect, got any) {
	t.Helper()
	equalGo(expect, got).prefix("%s == %s", Quote(ast.Assignment(0, 1)), Quote(ast.Assignment(0, 2))).false(t)
}

func equalGo(expect, got any) res {
	e := fmt.Sprintf("%#v", expect)
	g := fmt.Sprintf("%#v", got)
	if e == g {
		return res{true, e}
	} else {
		return res{false, fmt.Sprintf("%s != %s", Quote(e), Quote(g))}
	}
}

// compare using JSON
func EqualsJSON(t *testing.T, expect, got any) {
	t.Helper()
	equalJSON(expect, got).prefix("%s == %s", ast.Assignment(0, 1), ast.Assignment(0, 2)).true(t)
}

// compare using JSON
func NotEqualsJSON(t *testing.T, expect, got any) {
	t.Helper()
	equalJSON(expect, got).prefix("%s == %s", ast.Assignment(0, 1), ast.Assignment(0, 2)).false(t)
}

func equalJSON(expect, got any) res {
	e := jsonish(expect)
	g := jsonish(got)
	if e == g {
		return res{true, e}
	} else {
		return res{false, fmt.Sprintf("%s != %s", e, g)}
	}
}
