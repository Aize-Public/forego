package enc_test

import (
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

type X struct {
	Type string
	Path []string
}

func (this X) String() string {
	s := this.Type
	sep := ":"
	for _, p := range this.Path {
		s += sep + url.QueryEscape(p)
		sep = "/"
	}
	return s
}

var xRE = regexp.MustCompile(`^([a-z]+):([a-z]+(?:\/[a-z]+)*)$`)

func (this *X) Parse(c ctx.C, s string) error {
	out := xRE.FindStringSubmatch(s)
	if len(out) == 0 {
		return ctx.NewErrorf(c, "invalid X: %q", s)
	}
	log.Warnf(c, "out: %#v", out)
	this.Type = out[1]
	this.Path = []string{}
	for _, p := range strings.Split(out[2], "/") {
		this.Path = append(this.Path, url.QueryEscape(p))
	}
	return nil
}

var _ enc.Marshaler = &X{}
var _ enc.Unmarshaler = &X{}

func (this X) MarshalNode(c ctx.C) (enc.Node, error) {
	return enc.String(this.String()), nil
}

func (this *X) UnmarshalNode(c ctx.C, n enc.Node) error {
	switch n := n.(type) {
	case enc.String:
		return this.Parse(c, string(n))
	default:
		return ctx.NewErrorf(c, "expected string, got %T", n)
	}
}

func TestReadme(t *testing.T) {
	c := test.Context(t)

	a := X{
		Type: "xxx",
		Path: strings.Split("foo/bar/cuz", "/"),
	}
	t.Logf("a: %+v", a)

	n := enc.MustMarshal(c, a)
	t.Logf("n: %+v", n)

	var b X
	enc.MustUnmarshal(c, n, &b)
	t.Logf("b: %+v", b)

	test.EqualsGo(t, a, b)
}
