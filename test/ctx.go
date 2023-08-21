package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"golang.org/x/sys/unix"
)

// this ugly thing is true if the output goes to a console, false if the output is piped somewhere
var isTerminal = func() bool {
	_, err := unix.IoctlGetTermios(int(os.Stdout.Fd()), unix.TCGETS)
	return err == nil
}()

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
