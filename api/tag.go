package api

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/Aize-Public/forego/ctx"
)

type tag struct {
	name     string
	auth     bool
	in       bool
	out      bool
	required bool
	url      *url.URL
}

func tagName(c ctx.C, f reflect.StructField) string {
	var name string
	enc := f.Tag.Get("enc")
	json := f.Tag.Get("json")

	if enc != "" {
		name, _, _ = strings.Cut(enc, ",") // honor enc first
	} else if json != "" {
		name, _, _ = strings.Cut(json, ",") // then json
	}
	if name == "" { // if still no name, use field name
		name = f.Name
	}

	return name
}

func parseTags(c ctx.C, f reflect.StructField) (tag tag, err error) {
	parts := strings.Split(f.Tag.Get("api"), ",")
	tag.name = tagName(c, f)
	if tag.name == "" {
		tag.name = f.Name // fallback to field name
	}
	if u := f.Tag.Get("url"); u != "" {
		var err error
		tag.url, err = url.Parse(u)
		if err != nil {
			return tag, ctx.NewErrorf(c, "can't parse url: %q", u)
		}
	}

	for _, p := range parts {
		//log.Debugf(c, "%s %s", f.Name, p)
		switch p {
		case "in":
			tag.in = true
		case "out":
			tag.out = true
		case "both":
			tag.in = true
			tag.out = true
		case "required":
			tag.required = true
		case "auth":
			tag.auth = true
		case "":
		default:
			return tag, fmt.Errorf("invalid tag: %q", p)
		}
	}
	return tag, nil
}
