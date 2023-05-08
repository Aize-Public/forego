package test

import (
	"strings"
	"testing"
)

func ContainsJSON(t *testing.T, obj any, pattern string) {
	t.Helper()
	s := jsonish(obj)
	if strings.Contains(s, pattern) {
		t.Logf("contains %q: %q", pattern, obj)
	} else {
		t.Fatalf("FAIL: contains %q: %q", pattern, obj)
	}
}
