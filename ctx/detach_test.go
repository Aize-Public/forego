package ctx_test

import (
	"testing"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/test"
)

func TestDetach(t *testing.T) {
	c := test.Context(t)

	c1, cf1 := ctx.WithTimeout(c, time.Second)
	c1 = ctx.WithTag(c1, "tag", "1")

	c2 := ctx.Detach(c1)

	test.Assert(t, test.ExpectedError == ctx.RangeTag(c2, func(k string, _ ctx.JSON) error {
		if k == "tag" {
			return test.ExpectedError
		}
		return nil
	}))

	cf1() // cancel won't pass down...
	select {
	case <-c2.Done():
		test.Fail(t, "should not have been cancelled")
	default:
		test.Nil(t, c2.Err())
		deadline, exists := c2.Deadline()
		test.Assert(t, !exists)
		test.Assert(t, time.Time{} == deadline)
	}
}
