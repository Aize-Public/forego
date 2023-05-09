package test

import (
	"fmt"
	"testing"
)

// compare using "%#v"
func EqualsGo(t *testing.T, expect, got any) {
	t.Helper()
	equalGo(expect, got).true(t)
}

// compare using "%#v"
func NotEqualsGo(t *testing.T, expect, got any) {
	t.Helper()
	equalGo(expect, got).false(t)
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
	equalJSON(expect, got).true(t)
}

// compare using JSON
func NotEqualsJSON(t *testing.T, expect, got any) {
	t.Helper()
	equalJSON(expect, got).false(t)
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
