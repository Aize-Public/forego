package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Aize-Public/forego/ctx"
)

func NotContains(t *testing.T, str, pattern string) {
	t.Helper()
	contains(str, pattern).false(t)
}

func Contains(t *testing.T, str, pattern string) {
	t.Helper()
	contains(str, pattern).true(t)
}

func ContainsGo(t *testing.T, obj any, pattern string) {
	t.Helper()
	contains(fmt.Sprintf("%#v", obj), pattern).true(t)
}

// check if the json of obj contains pattern
func ContainsJSON(c ctx.C, obj any, pattern string) {
	t := ExtractTester(c)
	t.Helper()
	s := jsonish(c, obj)
	contains(s, pattern).true(t)
}

// check if the json of obj does NOT contains pattern
func NotContainsJSON(c ctx.C, obj any, pattern string) {
	t := ExtractTester(c)
	t.Helper()
	s := jsonish(c, obj)
	contains(s, pattern).false(t)
}

func contains(s string, pattern string) res {
	if strings.Contains(s, pattern) {
		return res{true, Quote(s)}
	} else {
		return res{false, fmt.Sprintf("%s not in %s", Quote(pattern), Quote(s))}
	}
}
