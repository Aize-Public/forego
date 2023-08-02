package prom_test

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/test"
	"github.com/Aize-Public/forego/utils/prom"
)

func TestHttp(t *testing.T) {
	c := test.Context(t)

	m := prom.Histogram{
		Desc:   "Foo Bar example",
		Labels: []string{"op", "loc"},
		Buckets: []float64{
			0.001, 0.002, 0.005,
			0.01, 0.02, 0.05,
			0.1, 0.2, 0.5,
			1, 2, 5,
		},
	}
	m.Observe(0.123, "foo", "bar")
	buf := &bytes.Buffer{}
	m.Print("foo_bar", buf)

	test.NoError(t, validHttpReponse(c, buf))
}

var metricRE = regexp.MustCompile(`^\w+{(\w+="[^"]+")(,\w+="[^"]+")*}\s\d+(\.\d+)?(\s\d+)?`)

func validHttpReponse(c ctx.C, buf *bytes.Buffer) error {
	for buf.Len() > 0 {
		ln, err := buf.ReadString('\n')
		if err != nil && err != io.EOF {
			return ctx.WrapError(c, err)
		}
		ln = strings.TrimRight(ln, " \r\n\t")
		if ln == "" {
			continue
		}
		switch ln[0] {
		case '#':
			log.Debugf(c, "HEAD %s", ln)
			// TODO
			//parts := strings.Split(ln, " ")
		default:
			if !metricRE.MatchString(ln) {
				return ctx.NewErrorf(c, "invalid line: %s", ln)
			}
		}
	}
	return nil
}
