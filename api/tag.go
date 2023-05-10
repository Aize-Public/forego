package api

import (
	"fmt"
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
	tag.in = true
	tag.out = true
	for _, p := range parts {
		//log.Debugf(c, "%s %s", f.Name, p)
		switch p {
		case "in":
			tag.out = false
		case "out":
			tag.in = false
		case "required":
			tag.required = true
		case "auth":
			tag.out = false
			tag.in = false
			tag.auth = true
		default:
			return tag, fmt.Errorf("invalid tag: %q", p)
		}
	}
	return tag, nil
}
