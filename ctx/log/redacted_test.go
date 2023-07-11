package log_test

import (
	"testing"

	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/test"
)

func TestRedacted(t *testing.T) {
	c := test.Context(t)
	var lines []log.Line
	k := log.WithLogger(c, func(line log.Line) {
		lines = append(lines, line)
	})
	s := log.RedactedString("foo")
	log.Debugf(k, "redacted %s string", s)
	test.Assert(t, len(lines) == 1)
	test.NotContainsJSON(t, lines[0], "foo")
}
