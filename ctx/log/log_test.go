package log_test

import (
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
	log.Debugf(c, "num: %d", 42)
	t.Logf("logs: %+v", logs)
	test.Assert(t, len(logs) == 1)
	test.ContainsJSON(t, logs[0], "num: 42")
}
