package log_test

import (
	"io"
	"strings"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/test"
)

func TestLogger(t *testing.T) {
	logs := []log.Line{}
	c := log.WithLogger(ctx.TODO(), func(at log.Line) {
		logs = append(logs, at)
	})
	c = ctx.WithTag(c, "foo", "bar")
	{
		log.Debugf(c, "num: %d", 42)
		t.Logf("logs: %+v", logs)
		test.Assert(t, len(logs) == 1)
		test.ContainsJSON(t, logs[0], "num: 42")
		test.ContainsJSON(t, logs[0].Tags["foo"], "bar")
	}
	{
		err := ctx.WrapError(c, io.EOF)
		log.Debugf(c, "err: %v", err)
		test.Assert(t, len(logs) == 2)
		errLog := logs[1]
		t.Logf("err: %+v", errLog)
		test.NotEmpty(t, errLog.Time)
		test.NotEmpty(t, errLog.Src)
		test.EqualsJSON(t, "debug", errLog.Level)
		test.ContainsJSON(t, errLog.Message, "EOF")
		test.NotEmpty(t, errLog.Tags)
		test.NotNil(t, errLog.Tags["error"])
		test.ContainsJSON(t, errLog.Tags["error"].String(), "EOF")
	}
	{
		err1 := ctx.WrapError(c, io.EOF)
		err2 := io.EOF
		log.Debugf(c, "err1: %v, err2: %v", err1, err2)
		test.Assert(t, len(logs) == 3)
		errLog := logs[2]
		t.Logf("err: %+v", errLog)
		test.NotEmpty(t, errLog.Tags)
		test.NotNil(t, errLog.Tags["error"])
		test.Assert(t, strings.HasPrefix(errLog.Tags["error"].String(), "["))
		test.Assert(t, strings.Count(errLog.Tags["error"].String(), "EOF") == 2)
	}
}
