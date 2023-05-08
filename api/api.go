package api

import (
	"context"
	"reflect"

	"github.com/Aize-Public/forego/ctx"
)

type Handler[T any] struct {
	typ reflect.Type

	in  []field
	out []field
}

func NewHandler[T any](c context.Context, init T) (Handler[T], error) {
	initV := reflect.ValueOf(init)
	this := Handler[T]{
		typ: initV.Type(),
	}

	for i := 0; i < this.typ.NumField(); i++ {
		ft := this.typ.Field(i)
		tag, err := parseTags(c, ft)
		if err != nil {
			return this, err
		}
		f := field{i, tag}
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

func (this Handler[T]) RequestOut(c ctx.C, obj T, data Data) error {
	v := reflect.ValueOf(obj)
	for _, f := range this.out {
		fv := v.Field(f.i)
		err := data.Marshal(c, f.tag.name, fv)
		if err != nil {
			return ctx.NewErrorf(c, "can't requestOut %T.%s: %w", obj, f.tag.name, err)
		}
	}
	return nil
}

func (this Handler[T]) RequestIn(c context.Context, data Data) (out T, err error) {
	ptrV := reflect.New(this.typ)
	v := ptrV.Elem()
	for _, f := range this.in {
		fv := v.Field(f.i)
		err := data.Unmarshal(c, f.tag.name, fv)
		if err != nil {
			return out, ctx.NewErrorf(c, "can't RequestInt %T.%s: %w", out, f.tag.name, err)
		}
	}
	return v.Interface().(T), nil
}

func (this Handler[T]) ResponseOut(c context.Context, obj T, data Data) (err error) {
	v := reflect.ValueOf(obj)
	for _, f := range this.out {
		fv := v.Field(f.i)
		err := data.Marshal(c, f.tag.name, fv)
		if err != nil {
			return ctx.NewErrorf(c, "can't ResponseOut %T.%s: %w", obj, f.tag.name, err)
		}
	}
	return nil
}
