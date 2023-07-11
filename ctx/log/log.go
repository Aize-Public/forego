package log

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Aize-Public/forego/ctx"
)

type Line struct {
	Level   string    `json:"level,omitempty"`
	Src     string    `json:"src,omitempty"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
	Tags    Tags      `json:"tags,omitempty"`
}

type Tags map[string]ctx.JSON

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

// return a new context with a custom logger attached to it
func WithLogger(c ctx.C, logger func(Line)) ctx.C {
	return context.WithValue(c, loggerKey{}, loggerValue{func() {}, logger})
}

// same as WithLogger(), but it has an extra helper function mostly used for testing
func WithLoggerAndHelper(c ctx.C, logger func(Line), helper func()) ctx.C {
	return context.WithValue(c, loggerKey{}, loggerValue{helper, logger})
}

var defLogger = func(at Line) {
	j, _ := json.Marshal(at)
	_, _ = fmt.Printf("%s\n", j)
}

func getLogger(c ctx.C) loggerValue {
	if c == nil {
		return loggerValue{func() {}, defLogger}
	}
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
		Src:   caller(1),
		Level: "error",
	}.formatf(c, f, args...))
}

func Warnf(c ctx.C, f string, args ...any) {
	l := getLogger(c)
	l.helper()
	l.log(Line{
		Src:   caller(1),
		Level: "warn",
	}.formatf(c, f, args...))
}

func Infof(c ctx.C, f string, args ...any) {
	l := getLogger(c)
	l.helper()
	l.log(Line{
		Src:   caller(1),
		Level: "info",
	}.formatf(c, f, args...))
}

func Debugf(c ctx.C, f string, args ...any) {
	l := getLogger(c)
	l.helper()
	l.log(Line{
		Src:   caller(1),
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
	at.Tags = tags(c)
LOOP:
	for i := 0; i < len(args); i++ {
		switch arg := args[i].(type) {
		case Loggable:
			// if it implements Loggable, replace it with the return value and allow manipulation of tags
			args[i] = arg.LogAs(&at.Tags)
		case error:
			// since errors can be wrapped, we need to unrwap them into ctx.Error to find the stack trace
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
	at.Message = fmt.Sprintf(f, args...)
	return at
}

type Loggable interface {
	LogAs(*Tags) any
}
