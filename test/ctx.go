package test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

func Context(t *testing.T) ctx.C {
	t.Helper()
	c := context.Background()
	c = ctx.WithTag(c, "test", t.Name())
	c = log.WithTester(c, t)
	c = log.WithLogger(c, slog.New(&TestLogger{Logger: log.NewDefaultLogger(os.Stdout)}))

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

type TestLogger struct {
	Logger *slog.Logger
}

var _ slog.Handler = &TestLogger{}

func (this *TestLogger) Handle(c context.Context, record slog.Record) error {
	if !isTerminal { // TODO(oha) allow for an env variable to override
		return this.Logger.Handler().Handle(c, record)
	} else {
		t := log.GetTester(c)
		t.Helper()
		t.Logf("%s: %s", record.Level, record.Message)
	}
	return nil
}

func (this *TestLogger) Enabled(c context.Context, level slog.Level) bool {
	return true
}

func (this *TestLogger) WithAttrs(attrs []slog.Attr) slog.Handler {
	return this
}

func (this *TestLogger) WithGroup(name string) slog.Handler {
	return this
}
