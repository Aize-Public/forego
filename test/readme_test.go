package test_test

import (
	"testing"

	"github.com/Aize-Public/forego/test"
)

func foo() error       { return nil }
func add(a, b int) int { return a + b }

// only needed to make sure readme examples are correct
func TestReadme(t *testing.T) {
	t.SkipNow()

	err := foo()
	test.NoError(t, err)

	test.EqualsGo(t, 2*2, add(2, 2))
	test.EqualsGo(t, 2*2, add(2, 3))
}

func TestReadAssert(t *testing.T) {
	t.SkipNow()

	test.Assert(t, 3 > 2)
	test.Assert(t, 3 <= 2)
}
