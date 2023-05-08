package test

import "testing"

func TestContains(t *testing.T) {
	contains("foobar", "oo").ok(t)
	contains("foobar", "").ok(t)
	contains("foobar", "cuz").fail(t)
}
