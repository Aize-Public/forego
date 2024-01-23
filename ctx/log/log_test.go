package log_test

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

type expectedLogStruct struct {
	Level string    `json:"level"`
	Msg   string    `json:"message"`
	Time  time.Time `json:"time"`
	Src   string    `json:"src"`
	Tags  log.Tags  `json:"tags"`
}

type loggableArg struct {
	value       string
	replaceTags log.Tags
}

var _ log.Loggable = loggableArg{}

func (this loggableArg) LogAs(tags *log.Tags) any {
	for k, v := range this.replaceTags {
		(*tags)[k] = v
	}
	return this.value
}

func TestLogger(t *testing.T) {
	c := test.Context(t)

	// Add some tags
	c = ctx.WithTag(c, "a", "string")
	c = ctx.WithTag(c, "b", 42)
	c = ctx.WithTag(c, "c", map[string]bool{"1": true, "2": true, "3": false})
	c = ctx.WithTag(c, "d", []int{1, 2, 3})
	expectedTags := []byte(`{"a":"string","b":42,"c":{"1":true,"2":true,"3":false},"d":[1,2,3],"test":"TestLogger"}`)

	// Add logger with custom buffer
	buf := &bytes.Buffer{}
	c = log.WithLogger(c, log.NewDefaultLogger(buf))

	verify := func(c ctx.C, expectedLevel, expectedMsg string, emptySrc bool) expectedLogStruct {
		t.Helper()
		defer buf.Reset()
		t.Logf("TESTING JSON LOG LINE: %s", buf.String())

		var m map[string]any
		test.NoError(t, enc.UnmarshalJSON(c, buf.Bytes(), &m))
		if emptySrc {
			test.EqualsGo(t, 4, len(m)) // check for unexpected fields
		} else {
			test.EqualsGo(t, 5, len(m)) // check for unexpected fields
		}
		var l expectedLogStruct
		test.NoError(t, enc.UnmarshalJSON(c, buf.Bytes(), &l))

		if emptySrc {
			test.Empty(t, l.Src)
		} else {
			test.NotEmpty(t, l.Src)
			_, filepath, _, _ := runtime.Caller(1)
			test.Assert(t, strings.HasPrefix(string(l.Src), filepath))
		}
		test.EqualsGo(t, expectedLevel, l.Level)
		test.EqualsGo(t, expectedMsg, l.Msg)
		test.NotEmpty(t, l.Time)
		tΔ := time.Since(l.Time)
		test.Assert(t, tΔ > 0)
		test.Assert(t, tΔ < time.Minute)
		return l
	}

	{
		log.Debugf(c, "Testing testing %d", 123)
		l := verify(c, "debug", "Testing testing 123", false)
		test.EqualsJSON(c, expectedTags, l.Tags)
	}
	{
		log.Infof(c, "Testing testing %d", 123)
		l := verify(c, "info", "Testing testing 123", false)
		test.EqualsJSON(c, expectedTags, l.Tags)
	}
	{
		log.Warnf(c, "Testing testing %d", 123)
		l := verify(c, "warn", "Testing testing 123", false)
		test.EqualsJSON(c, expectedTags, l.Tags)
	}
	{
		log.Errorf(c, "Testing testing %d%s%d", 1, "2", 3)
		l := verify(c, "error", "Testing testing 123", false)
		test.EqualsJSON(c, expectedTags, l.Tags)
	}
	{
		_, filepath, line, _ := runtime.Caller(0)
		log.Customf(c, slog.LevelDebug, fmt.Sprintf("%s:%d", filepath, line), "Testing testing %d", 123)
		l := verify(c, "debug", "Testing testing 123", false)
		test.EqualsJSON(c, expectedTags, l.Tags)
	}
	{ // Custom log level and no src
		log.Customf(c, slog.Level(int(slog.LevelError)+42), "", "Testing testing %d", 123)
		l := verify(c, "error+42", "Testing testing 123", true)
		test.EqualsJSON(c, expectedTags, l.Tags)
	}
	{ // Single wrapped error
		mockErr := ctx.WrapError(c, io.EOF)
		log.Errorf(c, "Testing error: %v", mockErr)
		l := verify(c, "error", "Testing error: EOF", false)
		test.NotEmpty(t, l.Tags["error"])

		var errs []map[string]any
		test.NoError(t, enc.UnmarshalJSON(c, []byte(l.Tags["error"].String()), &errs))
		test.EqualsGo(t, 1, len(errs))
		test.EqualsJSON(c, "EOF", errs[0]["error"])
		test.NotEmpty(t, errs[0]["stack"])
		test.NotEmpty(t, errs[0]["tags"])
		test.EqualsJSON(c, expectedTags, errs[0]["tags"])
	}
	{ // Multiple wrapped errors
		mockErr := ctx.WrapError(c, io.EOF)
		log.Errorf(c, "Testing error: err1=%v %v", mockErr, ctx.NewErrorf(c, "err2=%w", mockErr))
		l := verify(c, "error", "Testing error: err1=EOF err2=EOF", false)
		test.NotEmpty(t, l.Tags["error"])

		var errs []map[string]any
		test.NoError(t, enc.UnmarshalJSON(c, []byte(l.Tags["error"].String()), &errs))
		test.EqualsGo(t, 2, len(errs))
		test.EqualsJSON(c, "EOF", errs[0]["error"])
		test.EqualsJSON(c, "err2=EOF", errs[1]["error"])
		test.NotEmpty(t, errs[0]["stack"])
		test.NotEmpty(t, errs[1]["stack"])
		test.NotEmpty(t, errs[0]["tags"])
		test.NotEmpty(t, errs[1]["tags"])
		test.EqualsJSON(c, expectedTags, errs[0]["tags"])
		test.EqualsJSON(c, expectedTags, errs[1]["tags"])
	}
	{ // Loggable with no rewrite of tags
		arg := loggableArg{value: "Loggable arg", replaceTags: log.Tags{}}
		log.Infof(c, "Testing loggable: %v", arg)
		l := verify(c, "info", "Testing loggable: Loggable arg", false)
		test.EqualsJSON(c, expectedTags, l.Tags)
	}
	{ // Loggable with rewrite of tags
		arg := loggableArg{value: "Loggable arg", replaceTags: log.Tags{
			"a": ctx.JSON(`42`),
			"b": ctx.JSON(`"b"`),
			"c": ctx.JSON(`["1","2","3"]`),
			"d": ctx.JSON(`{"1":"yes","2":"no","3":"maybe"}`),
		}}
		log.Infof(c, "Testing loggable: %v", arg)
		l := verify(c, "info", "Testing loggable: Loggable arg", false)
		modifiedTags := []byte(`{"a":42,"b":"b","c":["1","2","3"],"d":{"1":"yes","2":"no","3":"maybe"},"test":"TestLogger"}`)
		test.EqualsJSON(c, modifiedTags, l.Tags)
	}
	{ // Loggables with multiple rewrites of the same tag, expecting last arg to win
		arg1 := loggableArg{value: "Arg1", replaceTags: log.Tags{
			"d": ctx.JSON(`{"1": 1}`),
		}}
		arg2 := loggableArg{value: "Arg2", replaceTags: log.Tags{
			"d": ctx.JSON(`2`),
		}}
		arg3 := loggableArg{value: "Arg3", replaceTags: log.Tags{
			"d": ctx.JSON(`3`),
		}}
		log.Infof(c, "Testing loggable: %v, %v, %v", arg1, arg2, arg3)
		l := verify(c, "info", "Testing loggable: Arg1, Arg2, Arg3", false)
		modifiedTags := []byte(`{"d":3,"a":"string","b":42,"c":{"1":true,"2":true,"3":false},"test":"TestLogger"}`)
		test.EqualsJSON(c, modifiedTags, l.Tags)
	}
	{ // Change the minimum log level
		log.SetDefaultLoggerLevel(slog.LevelInfo)
		log.Infof(c, "Testing testing %d", 123)
		log.Debugf(c, "Testing testing %d", 123) // this should not be logged now
		l := verify(c, "info", "Testing testing 123", false)
		test.EqualsJSON(c, expectedTags, l.Tags)
	}
}

func TestDefaultLogger(t *testing.T) {
	c, cf := ctx.Background()
	defer cf(nil)
	// now this goes to stdout, so it's not easy to verify the output, but at least we'll know if it panics
	log.Debugf(c, "Testing default logging %d", 42)
}

func BenchmarkLoggerDiscard(b *testing.B) {
	c, cf := ctx.Background()
	defer cf(nil)

	c = log.WithLogger(c, log.NewDefaultLogger(io.Discard))

	c = ctx.WithTag(c, "a", "string")
	c = ctx.WithTag(c, "b", 42)
	c = ctx.WithTag(c, "c", map[string]bool{"1": true, "2": true, "3": false})
	c = ctx.WithTag(c, "d", []int{1, 2, 3})

	for i := 0; i < b.N; i++ {
		log.Debugf(c, "Benching logger [%d]", i)
	}
}

func BenchmarkLoggerStdout(b *testing.B) {
	c, cf := ctx.Background()
	defer cf(nil)

	c = ctx.WithTag(c, "a", "string")
	c = ctx.WithTag(c, "b", 42)
	c = ctx.WithTag(c, "c", map[string]bool{"1": true, "2": true, "3": false})
	c = ctx.WithTag(c, "d", []int{1, 2, 3})

	for i := 0; i < b.N; i++ {
		log.Debugf(c, "Benching logger [%d]", i)
	}
}
