package log_test

import (
	"bytes"
	"testing"

	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/test"
)

func TestRedacted(t *testing.T) {
	c := test.Context(t)
	buf := &bytes.Buffer{}
	c = log.WithLogger(c, log.NewDefaultLogger(buf))
	s := log.RedactedString("foo")
	log.Debugf(c, "redacted %s string", s)
	test.NotContainsJSON(c, buf.String(), "foo")
	test.ContainsJSON(c, buf.String(), "***")
}
