package api

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type tag struct {
	name     string
	auth     bool
	in       bool
	out      bool
	required bool
}

func parseTags(c context.Context, f reflect.StructField) (tag tag, err error) {
	parts := strings.Split(f.Tag.Get("api"), ",")
	tag.name = parts[0]
	if tag.name == "" {
		tag.name = f.Name // fallback to field name
	}
	tag.in = true
	tag.out = true
	for _, p := range parts[1:] {
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
			// TODO HEAD
		default:
			return tag, fmt.Errorf("invalid tag: %q", p)
		}
	}
	return tag, nil
}
