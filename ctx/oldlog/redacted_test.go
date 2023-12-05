package oldlog_test

import (
	"testing"

	"github.com/Aize-Public/forego/ctx/oldlog"
	"github.com/Aize-Public/forego/test"
)

func TestRedacted(t *testing.T) {
	c := test.Context(t)
	var lines []oldlog.Line
	k := oldlog.WithLogger(c, func(line oldlog.Line) {
		lines = append(lines, line)
	})
	s := oldlog.RedactedString("foo")
	oldlog.Debugf(k, "redacted %s string", s)
	test.Assert(t, len(lines) == 1)
	test.NotContainsJSON(t, lines[0], "foo")
}
