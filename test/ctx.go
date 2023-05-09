package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

func C(t *testing.T) ctx.C {
	t.Helper()
	c := context.Background()
	c = log.WithLoggerAndHelper(c, func(ln log.Line) {
		if testing.Verbose() {
			fmt.Println(ln.JSON())
		} else {
			t.Helper()
			t.Logf("%s: %s", ln.Level, ln.Message)
			//t.Logf("%s %s: %s", ln.Level,
			//	filepath.Join(
			//		filepath.Base(filepath.Dir(ln.Src)),
			//		filepath.Base(ln.Src),
			//	), ln.Message)
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
