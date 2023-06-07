package api

import (
	"io"
	"reflect"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/enc"
)

// implementation for all the client/server request/responses using json
// Note(oha): the object is not goroutine safe, but it's not expected to be
type JSON struct {
	h enc.Handler
	j enc.JSON

	Data enc.Map
	UID  enc.Node
}

var _ ClientRequest = &JSON{}
var _ ServerRequest = &JSON{}
var _ ServerResponse = &JSON{}
var _ ClientResponse = &JSON{}

func (this JSON) String() string {
	j, _ := enc.MarshalJSON(ctx.TODO(), this)
	return string(j)
}

func (this *JSON) ReadFrom(c ctx.C, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return ctx.NewErrorf(c, "can't read api.JSON: %w", err)
	}
	if len(data) == 0 {
		this.Data = enc.Map{}
		return nil
	}
	n, err := this.j.Decode(c, data)
	if err != nil {
		return err
	}
	switch n := n.(type) {
	case enc.Map:
		this.Data = n
	default:
		return ctx.NewErrorf(c, "expected object, got %s", data)
	}
	return nil
}

func (this *JSON) Auth(c ctx.C, into reflect.Value, required bool) error {
	if (this.UID == nil || this.UID == enc.Nil{}) {
		if required {
			return ctx.NewErrorf(c, "Auth required")
		}
		return nil
	}
	return this.h.Unmarshal(c, this.UID, into.Addr().Interface())
}

func (this *JSON) Marshal(c ctx.C, name string, from reflect.Value) error {
	if this.Data == nil {
		this.Data = enc.Map{}
	}
	n, err := this.h.Marshal(c, from.Interface())
	if err != nil {
		return ctx.NewErrorf(c, "can't Marshal %q: %w", name, err)
	}
	this.Data[name] = n
	return nil
}

func (this *JSON) Unmarshal(c ctx.C, name string, into reflect.Value) error {
	if this.Data == nil {
		this.Data = enc.Map{}
	}
	n, ok := this.Data[name]
	if !ok {
		return nil
	}
	err := this.h.Unmarshal(c, n, into.Addr().Interface())
	//err := json.Unmarshal(j, into.Addr().Interface())
	if err != nil {
		return ctx.NewErrorf(c, "can't Unmarshal %q: %w", name, err)
	}
	return nil
}
