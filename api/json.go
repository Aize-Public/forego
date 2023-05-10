package api

import (
	"encoding/json"
	"reflect"

	"github.com/Aize-Public/forego/ctx"
)

// implementation for all the client/server request/responses using json
// Note(oha): the object is not goroutine safe, but it's not expected to be
type JSON struct {
	Data JSONMap
	UID  json.RawMessage // used by ServerRequest.Auth()
}

type JSONMap map[string]json.RawMessage

func (m JSONMap) String() string {
	j, _ := json.Marshal(m)
	return string(j)
}

var _ ClientRequest = &JSON{}
var _ ServerRequest = &JSON{}
var _ ServerResponse = &JSON{}
var _ ClientResponse = &JSON{}

func (this JSON) String() string {
	j, _ := json.Marshal(this)
	return string(j)
}

func (this *JSON) Auth(c ctx.C, into reflect.Value, required bool) error {
	if len(this.UID) == 0 || string(this.UID) == "null" {
		if required {
			return ctx.NewErrorf(c, "Auth required")
		}
		return nil
	}
	return json.Unmarshal(this.UID, into.Addr().Interface())
}

func (this *JSON) Marshal(c ctx.C, name string, from reflect.Value) error {
	if this.Data == nil {
		this.Data = map[string]json.RawMessage{}
	}
	j, err := json.Marshal(from.Interface())
	if err != nil {
		return ctx.NewErrorf(c, "can't Marshal %q: %w", name, err)
	}
	//log.Debugf(c, "json[%q]=%s", name, j)
	this.Data[name] = j
	return nil
}

func (this *JSON) Unmarshal(c ctx.C, name string, into reflect.Value) error {
	if this.Data == nil {
		this.Data = map[string]json.RawMessage{}
	}
	j, ok := this.Data[name]
	if !ok {
		return nil
	}
	err := json.Unmarshal(j, into.Addr().Interface())
	if err != nil {
		return ctx.NewErrorf(c, "can't Unmarshal %q: %w", name, err)
	}
	return nil
}
