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

	"github.com/Aize-Public/forego/ctx"
)

// Log arguments can implement this to rewrite themselves, or context tags, before being logged
type Loggable interface {
	LogAs(*Tags) any
}

type Tags map[string]ctx.JSON

// Converts the map of tags to a slog-friendly list
func (this *Tags) AsList() []any {
	res := make([]any, 0, len(*this)*2)
	for k, v := range *this {
		res = append(res, k)
		res = append(res, v)
	}
	return res
}

// Extracts any tags attached to the context
func ExtractTags(c ctx.C) Tags {
	out := Tags{}
	_ = ctx.RangeTag(c, func(k string, j ctx.JSON) error {
		out[k] = j
		return nil
	})
	return out
}

type LogFunc func(c ctx.C, level slog.Level, src, f string, args ...any)

var defaultLogFunc = wrapSlogLogger(NewDefaultSlogLogger(os.Stdout))

// This enables changing the minimum level of the default logger dynamically
var DefaultLoggerLevel = new(slog.LevelVar)

// Just a wrapper for the DefaultLoggerLevel.Set() method.
// This only applies to the default logger, unless you use the DefaultLoggerLevel variable also in your custom handler.
func SetDefaultLoggerLevel(level slog.Level) {
	DefaultLoggerLevel.Set(level)
}

// Returns a new context with the custom slog logger attached,
// with automatic handling of tags and Loggable arguments.
func WithSlogLogger(c ctx.C, l *slog.Logger) ctx.C {
	return WithLogFunc(c, wrapSlogLogger(l))
}

func wrapSlogLogger(l *slog.Logger) LogFunc {
	return func(c ctx.C, level slog.Level, src, f string, args ...any) {
		tags := ExtractTags(c)
		msg := FormatMsg(c, &tags, f, args...)
		if src != "" {
			l.LogAttrs(c, level, msg, slog.String("src", src), slog.Group("tags", tags.AsList()...))
		} else {
			l.LogAttrs(c, level, msg, slog.Group("tags", tags.AsList()...))
		}
	}
}

// Returns a new context with the custom LogFunc attached,
// but without handling of tags and Loggable arguments.
func WithLogFunc(c ctx.C, l LogFunc) ctx.C {
	conf := extractConfig(c)
	conf.logFunc = l
	return withConfig(c, conf)
}

// Returns a new context with the specified helper func attached,
// which is called automatically by functions through the log stack.
// Useful for tests if you do logging with the t.Logf func (wrapped in a LogFunc),
// and add the t.Helper func with this - ensuring that the relevant src line is logged.
func WithHelper(c ctx.C, h func()) ctx.C {
	conf := extractConfig(c)
	conf.helper = h
	return withConfig(c, conf)
}

// Creates a slog JSON logger with a certain default configuration, with the default minimum log level of debug
func NewDefaultSlogLogger(out io.Writer) *slog.Logger {
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

func Errorf(c ctx.C, f string, args ...any) {
	conf := extractConfig(c)
	helper(conf)()
	doLog(c, conf, slog.LevelError, caller(1), f, args...)
}

func Warnf(c ctx.C, f string, args ...any) {
	conf := extractConfig(c)
	helper(conf)()
	doLog(c, conf, slog.LevelWarn, caller(1), f, args...)
}

func Infof(c ctx.C, f string, args ...any) {
	conf := extractConfig(c)
	helper(conf)()
	doLog(c, conf, slog.LevelInfo, caller(1), f, args...)
}

func Debugf(c ctx.C, f string, args ...any) {
	conf := extractConfig(c)
	helper(conf)()
	doLog(c, conf, slog.LevelDebug, caller(1), f, args...)
}

// Log with a custom log level and src.
// To drop the src Attr entirely, for slog loggers, leave the string empty.
func Customf(c ctx.C, level slog.Level, src, f string, args ...any) {
	conf := extractConfig(c)
	helper(conf)()
	doLog(c, conf, level, src, f, args...)
}

func doLog(c ctx.C, conf *config, level slog.Level, src, f string, args ...any) {
	helper(conf)()

	logFunc := defaultLogFunc
	if conf.logFunc != nil {
		logFunc = conf.logFunc
	}
	logFunc(c, level, src, f, args...)
}

func helper(conf *config) func() {
	if conf.helper != nil {
		return conf.helper
	}
	return func() {}
}

type configKey struct{}

type config struct {
	logFunc LogFunc
	helper  func()
}

func withConfig(c ctx.C, conf *config) ctx.C {
	return context.WithValue(c, configKey{}, conf)
}

func extractConfig(c ctx.C) *config {
	if c != nil {
		if conf, ok := c.Value(configKey{}).(*config); ok {
			return conf
		}
	}
	return &config{}
}

func FormatMsg(c ctx.C, tags *Tags, f string, args ...any) string {
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
					m["tags"] = ExtractTags(err.C)
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
