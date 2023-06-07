package api

import (
	"net/url"
	"reflect"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

// each type of object T has its own handler
type Handler[T any] struct {
	typ  reflect.Type
	urls []*url.URL

	auth *field
	init map[int]reflect.Value
	in   []field
	out  []field
}

func NewServer[T any](c ctx.C, init T) (Server[T], error) {
	h, err := newHandler(c, init)
	return Server[T]{h}, err
}

func NewClient[T any](c ctx.C, obj T) (Client[T], error) {
	h, err := newHandler(c, obj)
	return Client[T]{h}, err
}

// we accept either Type or *Type
func newHandler[T any](c ctx.C, init T) (Handler[T], error) {
	log.Debugf(c, "NewHandler[%T]", init)
	initV := reflect.ValueOf(init)
	this := Handler[T]{
		init: map[int]reflect.Value{},
	}
	if initV.Kind() != reflect.Pointer {
		return this, ctx.NewErrorf(c, "expected *struct, got %T", init)
	}
	if initV.IsZero() {
		initV = reflect.New(reflect.TypeOf(init).Elem())
	}
	//log.Debugf(c, "initV: %+v", initV)
	initV = initV.Elem()
	this.typ = initV.Type()
	if this.typ.Kind() != reflect.Struct {
		return this, ctx.NewErrorf(c, "expected *struct, got %T", init)
	}

	// TODO(oha) use init as initializer
	for i := 0; i < this.typ.NumField(); i++ {
		ft := this.typ.Field(i)
		tag, err := parseTags(c, ft)
		if err != nil {
			return this, err
		}
		log.Debugf(c, "%v.%s %+v", this.typ, ft.Name, tag)
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
		if tag.url != nil {
			this.urls = append(this.urls, tag.url)
		}
		if !tag.in && !tag.out && tag.url == nil {
			v := initV.Field(f.i)
			if !v.IsZero() {
				this.init[f.i] = v
				log.Debugf(c, "init %v.%s = %#v", this.typ, f.tag.name, v)
			}
		}
	}

	return this, nil
}

func (this *Handler[T]) URL() *url.URL {
	if len(this.urls) > 0 {
		return this.urls[0]
	}
	return nil
}

func (this *Handler[T]) URLs() []*url.URL {
	return this.urls
}

type field struct {
	i   int
	tag tag
}

type Server[T any] struct {
	Handler[T]
}

type Client[T any] struct {
	Handler[T]
}

func (this Client[T]) Send(c ctx.C, obj T, data ClientRequest) error {
	v := reflect.ValueOf(obj).Elem()
	for _, f := range this.in {
		fv := v.Field(f.i)
		err := data.Marshal(c, f.tag.name, fv)
		if err != nil {
			return ctx.NewErrorf(c, "can't SendRequest %T.%s: %w", obj, f.tag.name, err)
		}
	}
	return nil
}

func (this Server[T]) Recv(c ctx.C, req ServerRequest) (T, error) {
	var zero T
	log.Debugf(c, "Server[%T].Recv(%+v)", zero, req)
	ptrV := reflect.New(this.typ)
	v := ptrV.Elem()
	for i, fv := range this.init {
		v.Field(i).Set(fv)
		//log.Debugf(c, "init %T.%v = %#v", this.typ, this.typ.Field(i).Name, fv)
	}
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
	return v.Addr().Interface().(T), nil
}

func (this Server[T]) Send(c ctx.C, obj T, res ServerResponse) (err error) {
	v := reflect.ValueOf(obj).Elem()
	for _, f := range this.out {
		fv := v.Field(f.i)
		err := res.Marshal(c, f.tag.name, fv)
		if err != nil {
			return ctx.NewErrorf(c, "can't SendResponse %T.%s: %w", obj, f.tag.name, err)
		}
	}
	return nil
}

func (this Client[T]) Recv(c ctx.C, res ClientResponse, into T) (err error) {
	v := reflect.ValueOf(into).Elem()
	for _, f := range this.out {
		fv := v.Field(f.i)
		err := res.Unmarshal(c, f.tag.name, fv)
		if err != nil {
			return ctx.NewErrorf(c, "can't RecvResponse %T.%s: %w", into, f.tag.name, err)
		}
	}
	return nil
}
