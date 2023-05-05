package test_test

import (
	"testing"

	"github.com/Aize-Public/forego/test"
)

func TestAssert(t *testing.T) {
	test.RunOk(t, "true", func(t test.T) {
		yes := true
		test.Assert(t, yes)
	})
	test.RunFail(t, "true", func(t test.T) {
		no := false
		test.Assert(t, no)
	})
}
