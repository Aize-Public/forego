package utils_test

import (
	"testing"

	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils"
)

func outer(t *testing.T) []string {
	return inner(t)
}

func inner(t *testing.T) []string {
	return utils.Stack(0, 10)
}

func TestStack(t *testing.T) {
	stack := outer(t)
	for i, ln := range stack {
		t.Logf("stack[%d]: %s", i, ln)
	}
	test.Assert(t, len(stack) > 2)
}
