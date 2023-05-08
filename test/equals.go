package test

import (
	"fmt"
	"testing"
)

func EqualsJSON(t *testing.T, expect, got any) {
	equalJSON(expect, got).ok(t)
}

func NotEqualsJSON(t *testing.T, expect, got any) {
	equalJSON(expect, got).fail(t)
}

func equalJSON(expect, got any) res {
	e := jsonish(expect)
	g := jsonish(got)
	if e == g {
		return res{true, fmt.Sprintf("expected: %s", e)}
	} else {
		return res{false, fmt.Sprintf("expected %s got %s", e, g)}
	}
}
