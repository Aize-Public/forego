package log

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/utils"
)

type Line struct {
	Level   string                     `json:"level"`
	Src     string                     `json:"src"`
	Time    time.Time                  `json:"time"`
	Message string                     `json:"message"`
	Tags    map[string]json.RawMessage `json:"tags"`
}

func tags(c ctx.C) map[string]json.RawMessage {
	out := map[string]json.RawMessage{}
	_ = ctx.RangeTag(c, func(k string, j []byte) error {
		out[k] = j
		return nil
	})
	return out
}

type loggerKey struct{}

func WithLogger(c ctx.C, logger func(Line)) ctx.C {
	return context.WithValue(c, loggerKey{}, logger)
}

var defLogger = func(at Line) {
	j, _ := json.Marshal(at)
	_, _ = fmt.Printf("%s\n", j)
}

func logger(c ctx.C) func(Line) {
	logger, _ := c.Value(loggerKey{}).(func(Line))
	if logger == nil {
		return defLogger
	}
	return logger
}

func Debugf(c ctx.C, f string, args ...any) {
	logger(c)(Line{
		Src:   utils.Caller(1).FileLine(),
		Level: "debug",
	}.Formatf(c, f, args...))
}

func (at Line) Formatf(c ctx.C, f string, args ...any) Line {
	if at.Time.IsZero() {
		at.Time = time.Now()
	}
	at.Message = fmt.Sprintf(f, args...)
	at.Tags = tags(c)
LOOP:
	for i := 0; i < len(args); i++ {
		switch arg := args[i].(type) {
		case error:
			m := map[string]any{
				"error": arg.Error(),
			}
			var err ctx.Err
			if errors.As(arg, &err) {
				m["stack"] = err.Stack
				if err.C != nil {
					m["tags"] = tags(err.C)
				}
			}
			at.Tags["error"], _ = json.Marshal(m)
			break LOOP
		}
	}
	return at
}
