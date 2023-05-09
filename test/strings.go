package test

import (
	"fmt"
	"strings"
	"testing"
)

// check if the json of obj contains pattern
func ContainsJSON(t *testing.T, obj any, pattern string) {
	t.Helper()
	s := jsonish(obj)
	contains(s, pattern).true(t)
}

// check if the json of obj does NOT contains pattern
func NotContainsJSON(t *testing.T, obj any, pattern string) {
	t.Helper()
	s := jsonish(obj)
	contains(s, pattern).false(t)
}

func contains(s string, pattern string) res {
	if strings.Contains(s, pattern) {
		return res{true, s}
	} else {
		return res{false, fmt.Sprintf("%q not in %q", pattern, s)}
	}
}
