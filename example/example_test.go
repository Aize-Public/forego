package example_test

import (
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/test"
)

func lib(c ctx.C) {
	log.Debugf(c, "foobar")
}

func TestAll(t *testing.T) {
	c := test.Context(t)
	t.Logf("before")
	lib(c)
	t.Logf("after")
	//test.Assert(t, false)
}
