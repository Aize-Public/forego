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
	Level   string              `json:"level,omitempty"`
	Src     string              `json:"src,omitempty"`
	Time    time.Time           `json:"time"`
	Message string              `json:"message"`
	Tags    map[string]ctx.JSON `json:"tags,omitempty"`
}

func (this Line) JSON() string {
	j, err := json.Marshal(this)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func tags(c ctx.C) map[string]ctx.JSON {
	out := map[string]ctx.JSON{}
	_ = ctx.RangeTag(c, func(k string, j ctx.JSON) error {
		out[k] = j
		return nil
	})
	return out
}

type loggerKey struct {
}

type loggerValue struct {
	helper func()
	log    func(Line)
}

func WithLogger(c ctx.C, logger func(Line)) ctx.C {
	return context.WithValue(c, loggerKey{}, loggerValue{func() {}, logger})
}

func WithLoggerAndHelper(c ctx.C, logger func(Line), helper func()) ctx.C {
	return context.WithValue(c, loggerKey{}, loggerValue{helper, logger})
}

var defLogger = func(at Line) {
	j, _ := json.Marshal(at)
	_, _ = fmt.Printf("%s\n", j)
}

func getLogger(c ctx.C) loggerValue {
	logger, ok := c.Value(loggerKey{}).(loggerValue)
	if !ok {
		return loggerValue{func() {}, defLogger}
	}
	return logger
}

func Errorf(c ctx.C, f string, args ...any) {
	l := getLogger(c)
	l.helper()
	l.log(Line{
		Src:   utils.Caller(1).FileLine(),
		Level: "error",
	}.formatf(c, f, args...))
}

func Warnf(c ctx.C, f string, args ...any) {
	l := getLogger(c)
	l.helper()
	l.log(Line{
		Src:   utils.Caller(1).FileLine(),
		Level: "warn",
	}.formatf(c, f, args...))
}

func Infof(c ctx.C, f string, args ...any) {
	l := getLogger(c)
	l.helper()
	l.log(Line{
		Src:   utils.Caller(1).FileLine(),
		Level: "info",
	}.formatf(c, f, args...))
}

func Debugf(c ctx.C, f string, args ...any) {
	l := getLogger(c)
	l.helper()
	l.log(Line{
		Src:   utils.Caller(1).FileLine(),
		Level: "debug",
	}.formatf(c, f, args...))
}

// TODO(oha) needed?
func (at Line) Log(c ctx.C) {
	l := getLogger(c)
	l.helper()
	l.log(at)
}

func (at Line) formatf(c ctx.C, f string, args ...any) Line {
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
			var err ctx.Error
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
