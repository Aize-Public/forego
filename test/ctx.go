package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

func Context(t *testing.T) ctx.C {
	t.Helper()
	c := context.Background()
	c = ctx.WithTag(c, "test", t.Name())
	c = log.WithLoggerAndHelper(c, func(ln log.Line) {
		if !isTerminal { // TODO(oha) allow for an env variable to override
			fmt.Println(ln.JSON())
		} else {
			t.Helper()
			t.Logf("%s: %s", ln.Level, ln.Message)
		}
	}, t.Helper)
	d, ok := t.Deadline()
	if ok {
		c, cf := context.WithDeadline(c, d)
		t.Cleanup(cf)
		return c
	} else {
		c, cf := context.WithCancel(c)
		t.Cleanup(cf)
		return c
	}
}
