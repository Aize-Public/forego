package test

import (
	"fmt"
	"testing"

	"github.com/Aize-Public/forego/utils/ast"
)

// compare using "%#v"
func EqualsGo(t *testing.T, expect, got any) {
	t.Helper()
	if equalGo(expect, got).succeed {
		OK(t, "%s == %s: %#v", ast.Assignment(0, 1), ast.Assignment(0, 2), expect)
	} else {
		Fail(t, "%#v == %#v", expect, got)
	}
}

// compare using "%#v"
func NotEqualsGo(t *testing.T, expect, got any) {
	t.Helper()
	if equalGo(expect, got).succeed {
		Fail(t, "%s != %s", ast.Assignment(0, 1), ast.Assignment(0, 2))
	} else {
		OK(t, "%s != %s", ast.Assignment(0, 1), ast.Assignment(0, 2))
	}
}

func equalGo(expect, got any) res {
	e := fmt.Sprintf("%#v", expect)
	g := fmt.Sprintf("%#v", got)
	if e == g {
		return res{true, e}
	} else {
		return res{false, fmt.Sprintf("%s != %s", e, g)}
	}
}

// compare using JSON
func EqualsJSON(t *testing.T, expect, got any) {
	t.Helper()
	if equalJSON(expect, got).succeed {
		OK(t, "%s == %s: %#v", ast.Assignment(0, 1), ast.Assignment(0, 2), expect)
	} else {
		Fail(t, "%#v == %#v", expect, got)
	}
}

// compare using JSON
func NotEqualsJSON(t *testing.T, expect, got any) {
	t.Helper()
	if equalJSON(expect, got).succeed {
		Fail(t, "%s != %s", ast.Assignment(0, 1), ast.Assignment(0, 2))
	} else {
		OK(t, "%s != %s", ast.Assignment(0, 1), ast.Assignment(0, 2))
	}
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
