package api

import (
	"encoding/json"

	"github.com/Aize-Public/forego/ctx"
)

type Data interface {
	Unmarshal(c ctx.C, name string, into any) error
	Marshal(c ctx.C, name string, from any) error
}

type JSON map[string]json.RawMessage

var _ Data = JSON{}

func (this JSON) String() string {
	j, _ := json.Marshal(this)
	return string(j)
}

func (this JSON) Marshal(c ctx.C, name string, from any) error {
	j, err := json.Marshal(from)
	if err != nil {
		return ctx.NewErrorf(c, "can't Marshal %q: %w", name, err)
	}
	this[name] = j
	return nil
}

func (this JSON) Unmarshal(c ctx.C, name string, into any) error {
	j, ok := this[name]
	if !ok {
		return nil
	}
	err := json.Unmarshal(j, into)
	if err != nil {
		return ctx.NewErrorf(c, "can't Unmarshal %q: %w", name, err)
	}
	return nil
}
