package test

import (
	"testing"
)

func EqualsJSON(t testing.TB, expect, got any) {
	t.Helper()
	e := jsonish(expect)
	g := jsonish(got)
	if e == g {
		t.Logf("%s", e)
	} else {
		t.Fatalf("FAIL: expected %s got %s", e, g)
	}
}
