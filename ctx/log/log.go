package log

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/Aize-Public/forego/ctx"
)

// Log arguments can implement this to rewrite themselves, or context tags, before being logged
type Loggable interface {
	LogAs(*Tags) any
}

type Tags map[string]ctx.JSON

func (this *Tags) AsList() []any {
	res := make([]any, 0, len(*this)*2)
	for k, v := range *this {
		res = append(res, k)
		res = append(res, v)
	}
	return res
}

var defaultLogger = NewDefaultLogger(os.Stdout)

// This enables changing the minimum level of the default logger dynamically
var DefaultLoggerLevel = new(slog.LevelVar)

// Just a wrapper for the DefaultLoggerLevel.Set() method.
// This only applies to the default logger, unless you use the DefaultLoggerLevel variable also in your custom handler.
func SetDefaultLoggerLevel(level slog.Level) {
	DefaultLoggerLevel.Set(level)
}

// Creates a slog JSON logger with a certain default configuration, with the default minimum log level of debug
func NewDefaultLogger(out io.Writer) *slog.Logger {
	DefaultLoggerLevel.Set(slog.LevelDebug)
	return slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: DefaultLoggerLevel, ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.LevelKey:
			level := a.Value.Any().(slog.Level)
			return slog.String("level", strings.ToLower(level.String()))
		case slog.MessageKey:
			a.Key = "message"
		}
		return a
	}}))
}

// Returns a new context with the custom logger attached
func WithLogger(c ctx.C, logger *slog.Logger) ctx.C {
	conf := getConfig(c)
	conf.l = logger
	return withConfig(c, conf)
}

// If a valid logger is attached to the context, it's returned together with a status of `true`.
// If no valid logger is set on the context, it returns the default JSON logger instead, and `false`.
func GetLogger(c ctx.C) (*slog.Logger, bool) {
	conf := getConfig(c)
	if conf.l != nil {
		return conf.l, true
	}
	return defaultLogger, false
}

// Returns a new context with the testing object attached.
// This will in turn cause the log functions to call t.Helper() automatically
func WithTester(c ctx.C, t *testing.T) ctx.C {
	conf := getConfig(c)
	conf.t = t
	return withConfig(c, conf)
}

// Returns any testing object attached to the context, else nil
func GetTester(c ctx.C) *testing.T {
	return getConfig(c).t
}

func Errorf(c ctx.C, f string, args ...any) {
	conf := getConfig(c)
	helper(conf)()
	doLog(c, conf, slog.LevelError, caller(1), f, args...)
}

func Warnf(c ctx.C, f string, args ...any) {
	conf := getConfig(c)
	helper(conf)()
	doLog(c, conf, slog.LevelWarn, caller(1), f, args...)
}

func Infof(c ctx.C, f string, args ...any) {
	conf := getConfig(c)
	helper(conf)()
	doLog(c, conf, slog.LevelInfo, caller(1), f, args...)
}

func Debugf(c ctx.C, f string, args ...any) {
	conf := getConfig(c)
	helper(conf)()
	doLog(c, conf, slog.LevelDebug, caller(1), f, args...)
}

// Log with a custom log level and src. To drop the src Attr entirely, leave the string empty.
func Customf(c ctx.C, level slog.Level, src, f string, args ...any) {
	conf := getConfig(c)
	helper(conf)()
	doLog(c, conf, level, src, f, args...)
}

func doLog(c ctx.C, conf *config, level slog.Level, src, f string, args ...any) {
	helper(conf)()
	l := defaultLogger
	if conf.l != nil {
		l = conf.l
	}
	tags := extractTags(c)
	msg := formatMsg(c, &tags, f, args...)
	if src != "" {
		l.LogAttrs(c, level, msg, slog.String("src", src), slog.Group("tags", tags.AsList()...))
	} else {
		l.LogAttrs(c, level, msg, slog.Group("tags", tags.AsList()...))
	}
}

func helper(conf *config) func() {
	if conf.t != nil {
		return conf.t.Helper
	}
	return func() {}
}

type configKey struct{}

type config struct {
	l *slog.Logger
	t *testing.T
}

func withConfig(c ctx.C, conf *config) ctx.C {
	return context.WithValue(c, configKey{}, conf)
}

func getConfig(c ctx.C) *config {
	if c != nil {
		if conf, ok := c.Value(configKey{}).(*config); ok {
			return conf
		}
	}
	return &config{}
}

func extractTags(c ctx.C) Tags {
	out := Tags{}
	_ = ctx.RangeTag(c, func(k string, j ctx.JSON) error {
		out[k] = j
		return nil
	})
	return out
}

func formatMsg(c ctx.C, tags *Tags, f string, args ...any) string {
	errs := []map[string]any{}
	for i := 0; i < len(args); i++ {
		switch arg := args[i].(type) {
		case Loggable:
			// if it implements Loggable, replace it with the return value and allow manipulation of tags
			args[i] = arg.LogAs(tags)

		case error:
			m := map[string]any{"error": arg.Error()}
			var err ctx.Error
			if errors.As(arg, &err) {
				m["stack"] = err.Stack
				if err.C != nil {
					m["tags"] = extractTags(err.C)
				}
			}
			errs = append(errs, m)
		}
	}
	if len(errs) > 0 {
		if j, err := json.Marshal(errs); err == nil {
			(*tags)["error"] = ctx.JSON(j)
		}
	}
	return fmt.Sprintf(f, args...)
}

func caller(above int) string {
	_, file, line, _ := runtime.Caller(above + 1)
	return fmt.Sprintf("%s:%d", file, line)
}
