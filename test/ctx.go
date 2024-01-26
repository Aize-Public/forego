package test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

// Creates a new context for use in tests, with the testing object
// and a custom logger attached (which does logging with t.Logf).
// The context is cancelled automatically when the test ends.
func Context(t *testing.T) ctx.C {
	t.Helper()
	c := context.Background()
	c = ctx.WithTag(c, "test", t.Name())
	c = WithTester(c, t)

	if isTerminal { // TODO(oha) allow for an env variable to override
		c = log.WithHelper(c, t.Helper)
		c = log.WithLogFunc(c, func(c ctx.C, level slog.Level, src, f string, args ...any) {
			t.Helper()
			t.Logf("%s: %s", level, fmt.Sprintf(f, args...))
		})
	}

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

type testerKey struct{}

// Returns a new context with the testing object attached
func WithTester(c ctx.C, t testing.TB) ctx.C {
	return context.WithValue(c, testerKey{}, t)
}

// Returns any testing object attached to the context, else nil
func ExtractTester(c ctx.C) *testing.T {
	if c != nil {
		if t, ok := c.Value(testerKey{}).(*testing.T); ok {
			return t
		}
	}
	return nil
}
