package test

import "testing"

func TestContains(t *testing.T) {
	contains("foobar", "oo").true(t)
	contains("foobar", "").true(t)
	contains("foobar", "cuz").false(t)
}
