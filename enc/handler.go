package enc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

// generic object that do the Unmarshal()/Conflate()
type Handler struct {
	Factory map[reflect.Type]func(c ctx.C, n Node) (any, error)

	// called if a field is present in the NodeTree but there is no mapping on the object it's unmarshaled into
	UnhandledFields func(c ctx.C, path path, n Node) error

	Debugf func(c ctx.C, f string, args ...any)

	path path
}

func (this Handler) Append(p any) Handler {
	this.path = append(this.path, p)
	return this
}

var ineff = map[string]bool{}

func Unmarshal(c ctx.C, n Node, into any) error {
	return Handler{}.Unmarshal(c, n, into)
}

func (this Handler) Unmarshal(c ctx.C, n Node, into any) error {
	if n == nil {
		n = Nil{}
	}
	v := reflect.ValueOf(into).Elem()
	return this.Append(v.Type()).unmarshal(c, n, v)
}

func (this Handler) unmarshal(c ctx.C, from Node, v reflect.Value) error {
	c = ctx.WithTag(c, "path", this.path.String())
	if this.Debugf != nil {
		this.Debugf(c, "unmarshal( %v -> %v{%v} )", from, v.Type(), v)
	}
	if !v.CanSet() {
		return ctx.NewErrorf(c, "can't assign %v", v)
	}
	if v.Kind() == reflect.Pointer {
		switch from.(type) {
		case Nil:
			v.SetZero()
			return nil
		default:
			vv := reflect.New(v.Type().Elem())
			v.Set(vv)
			if this.Debugf != nil {
				this.Debugf(c, "pointer deref: %v -> %v", v.Type(), v.Elem().Type())
			}
			return this.unmarshal(c, from, v.Elem())
		}
	}
	switch into := v.Addr().Interface().(type) {
	case *time.Time:
		if this.Debugf != nil {
			this.Debugf(c, "is %T", into)
		}
		switch from := from.(type) {
		case String:
			var err error
			*into, err = time.Parse(time.RFC3339Nano, from.String())
			return err
		default:
			return ctx.NewErrorf(c, "can't expend time.Time from %T (%v)", from, from)
		}
	case Unmarshaler:
		if this.Debugf != nil {
			this.Debugf(c, "is %T", into)
		}
		return into.UnmarshalTree(c, from)
	case *json.RawMessage:
		if this.Debugf != nil {
			this.Debugf(c, "is %T", into)
		}
		*into = JSON{}.Encode(c, from)
		warnIneff(c, "Warn: inefficient json.RawMessage, use enc.Tree instead")
		return nil
	case json.Unmarshaler:
		if this.Debugf != nil {
			this.Debugf(c, "is %T", into)
		}
		warnIneff(c, "Warn: inefficient %T.UnmarshalJSON(): implement enc.Unmarshaler instead", into)
		j, _ := json.Marshal(from) // we must go back to the json
		return into.UnmarshalJSON(j)
	case *Node:
		if this.Debugf != nil {
			this.Debugf(c, "is %T", into)
		}
		// NOTE(oha) if a struct has a field of type enc.Node, we drop the data there (similarly to json.RawMessage)
		*into = from
		return nil
	}

	/*
		vv := v.Elem()
		switch from.(type) {
		case nil, Nil:
			vv.SetZero()
			return nil
		}
	*/

	if this.Factory != nil {
		f := this.Factory[v.Type()]
		if f != nil {
			log.Debugf(c, "factory for type %v", v.Type())
			obj, err := f(c, from)
			if err != nil {
				return ctx.NewErrorf(c, "factory error: %w", err)
			}
			s := reflect.ValueOf(obj)
			if !s.CanConvert(v.Elem().Type()) {
				return ctx.NewErrorf(c, "Factory %v returned %v which can't be converted to %v",
					v.Type().Elem(), s.Type(), v.Elem().Type())
			}
			log.Debugf(c, "Factory %v returned %v which can be converted to %v",
				v.Type().Elem(), s.Type(), v.Elem().Type())
			v.Elem().Set(s) // we set the value, not the pointer... because go.reflect
			return nil
		}
	}

	log.Debugf(c, "normal type: %v, use generic %T.unmarshalInto()", v.Type(), from)
	/*
		if v.Kind() == reflect.Pointer {
			// if we unmarshal into a pointer, we create the zero value and unmarshal into that instead
			// this allow to unmarshal into *bool or *string
			e := reflect.New(v.Type().Elem())
			log.Debugf(c, "OHA4: %v", e.Type())
			err := from.unmarshalInto(c, this, path, e.Elem())
			if err != nil {
				return err
			}
			v.Set(e)
			return nil
		}
	*/

	return from.unmarshalInto(c, this, v)
}

func warnIneff(c ctx.C, f string, args ...any) {
	msg := fmt.Sprintf(f, args...)
	if !ineff[msg] {
		ineff[msg] = true
		log.Warnf(c, f, args...)
	}
}

// transform an object into a enc.Node
func Marshal(c ctx.C, in any) (Node, error) {
	return Handler{}.Marshal(c, in)
}

// TODO(oha) no reason to have the Handler to Unmarshal()
func (this Handler) Marshal(c ctx.C, in any) (Node, error) {
	switch in := in.(type) {
	case nil:
		return Nil{}, nil
	case Marshaler:
		return in.MarshalTree(c)
	case json.Marshaler:
		j, err := in.MarshalJSON()
		if err != nil {
			return nil, err
		}
		return JSON{}.Decode(c, j)
	case Node:
		return in, nil
	}

	v := reflect.ValueOf(in)
	t := v.Type()
	switch t.Kind() {
	default:
		log.Warnf(c, "possible wrong fallback for type %T", in)
		return fromNative(in), nil
	case reflect.Bool:
		return Bool(v.Bool()), nil
	case reflect.Int:
		return Number(v.Int()), nil
	case reflect.String:
		return String(v.String()), nil
	case reflect.Pointer:
		if v.IsNil() {
			return Nil{}, nil
		}
		return this.Marshal(c, v.Elem().Interface())
	case reflect.Slice:
		if v.IsNil() {
			return Nil{}, nil
		}
		list := List{}
		for i := 0; i < v.Len(); i++ {
			ev := v.Index(i)
			e, err := this.Marshal(c, ev.Interface())
			if err != nil {
				return nil, err
			}
			list = append(list, e)
		}
		return list, nil
	case reflect.Map:
		if v.IsNil() {
			return Nil{}, nil
		}
		m := Map{}
		for _, kv := range v.MapKeys() {
			vv := v.MapIndex(kv)
			n, err := this.Marshal(c, vv.Interface())
			if err != nil {
				return nil, err
			}
			k := fmt.Sprint(kv.Interface()) // TODO(oha): is this ok?
			m[k] = n
		}
		return m, nil

	case reflect.Struct:
	}

	out := Pairs{}
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		if !ft.IsExported() {
			continue
		}
		tag := parseTag(ft)
		if tag.Skip {
			continue
		}
		fv := v.Field(i)
		fn, err := this.Marshal(c, fv.Interface())
		if err != nil {
			return nil, err
		}
		p := Pair{tag.Name, tag.JSON, fn}
		out = append(out, p)
	}
	return out, nil
}
