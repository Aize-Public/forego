package api

import (
	"reflect"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

// each type of object T has its own handler
type Handler[T any] struct {
	typ reflect.Type

	auth *field
	in   []field
	out  []field
}

func NewHandler[T any](c ctx.C, init T) (Handler[T], error) {
	log.Debugf(c, "NewHandler[%T]", init)
	initV := reflect.ValueOf(init)
	this := Handler[T]{
		typ: initV.Type(),
	}
	if this.typ.Kind() != reflect.Struct {
		return this, ctx.NewErrorf(c, "expected struct, got %T", init)
	}

	// TODO(oha) use init as initializer
	for i := 0; i < this.typ.NumField(); i++ {
		ft := this.typ.Field(i)
		tag, err := parseTags(c, ft)
		if err != nil {
			return this, err
		}
		f := field{i, tag}
		if tag.auth {
			if this.auth != nil {
				return this, ctx.NewErrorf(c, "only 1 auth field is supported: %T.%s", init, ft.Name)
			}
			this.auth = &f
		}
		if tag.in {
			this.in = append(this.in, f)
		}
		if tag.out {
			this.out = append(this.out, f)
		}
	}

	return this, nil
}

type field struct {
	i   int
	tag tag
}

func (this Handler[T]) Server() server[T] { return server[T]{this} }
func (this Handler[T]) Client() client[T] { return client[T]{this} }

type server[T any] struct {
	Handler[T]
}

type client[T any] struct {
	Handler[T]
}

func (this client[T]) Send(c ctx.C, obj T, data ClientRequest) error {
	v := reflect.ValueOf(obj)
	for _, f := range this.in {
		fv := v.Field(f.i)
		err := data.Marshal(c, f.tag.name, fv)
		if err != nil {
			return ctx.NewErrorf(c, "can't SendRequest %T.%s: %w", obj, f.tag.name, err)
		}
	}
	return nil
}

func (this server[T]) Recv(c ctx.C, req ServerRequest) (T, error) {
	var zero T
	ptrV := reflect.New(this.typ)
	v := ptrV.Elem()
	for _, f := range this.in {
		fv := v.Field(f.i)
		err := req.Unmarshal(c, f.tag.name, fv)
		if err != nil {
			return zero, ctx.NewErrorf(c, "can't RecvRequest %T.%s: %w", zero, f.tag.name, err)
		}
	}
	if this.auth != nil {
		fv := v.Field(this.auth.i)
		err := req.Auth(c, fv, this.auth.tag.required)
		if err != nil {
			return zero, ctx.NewErrorf(c, "can't RecvRequest %T Auth(): %w", zero, err)
		}
	}
	return v.Interface().(T), nil
}

func (this server[T]) Send(c ctx.C, obj T, res ServerResponse) (err error) {
	v := reflect.ValueOf(obj)
	for _, f := range this.out {
		fv := v.Field(f.i)
		err := res.Marshal(c, f.tag.name, fv)
		if err != nil {
			return ctx.NewErrorf(c, "can't SendResponse %T.%s: %w", obj, f.tag.name, err)
		}
	}
	return nil
}

func (this client[T]) Recv(c ctx.C, res ClientResponse, into T) (err error) {
	v := reflect.ValueOf(into)
	if v.Kind() != reflect.Pointer {
		return ctx.NewErrorf(c, "expected pointer, got %T", into)
	}
	for _, f := range this.out {
		fv := v.Field(f.i)
		err := res.Unmarshal(c, f.tag.name, fv)
		if err != nil {
			return ctx.NewErrorf(c, "can't RecvResponse %T.%s: %w", into, f.tag.name, err)
		}
	}
	return nil
}
