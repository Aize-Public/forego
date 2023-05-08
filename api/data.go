package api

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/Aize-Public/forego/ctx"
)

type Data interface {
	Unmarshal(c context.Context, name string, into reflect.Value) error
	Marshal(c context.Context, name string, from reflect.Value) error
}

type JSON map[string]json.RawMessage

var _ Data = JSON{}

func (this JSON) Marshal(c context.Context, name string, from reflect.Value) error {
	j, err := json.Marshal(from.Interface())
	if err != nil {
		return ctx.NewErrorf(c, "can't Marshal %q: %w", name, err)
	}
	this[name] = j
	return nil
}

func (this JSON) Unmarshal(c context.Context, name string, into reflect.Value) error {
	j, ok := this[name]
	if !ok {
		return nil
	}
	err := json.Unmarshal(j, into.Interface())
	if err != nil {
		return ctx.NewErrorf(c, "can't Unmarshal %q: %w", name, err)
	}
	return nil
}
