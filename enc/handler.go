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
	UnhandledFields func(c ctx.C, path []any, n Node) error

	Debugf func(c ctx.C, f string, args ...any)

	path path
}

func (this Handler) Append(p any) Handler {
	this.path = append(this.path, p)
	return this
}

var ineff = map[string]bool{}

func MustUnmarshal(c ctx.C, n Node, into any) {
	err := Handler{}.Unmarshal(c, n, into)
	if err != nil {
		panic(err)
	}
}

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
	//c = ctx.WithTag(c, "path", this.path.String()) // NOTE(oha): this is a bit slow because the json part
	//defer log.Debugf(c, "unmarshal( %T %+v => %v{%+v} )", from, from, v.Type(), v)
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
	case Unmarshaler:
		if this.Debugf != nil {
			this.Debugf(c, "is %T", into)
		}
		return into.UnmarshalNode(c, from)
	case *json.RawMessage:
		if this.Debugf != nil {
			this.Debugf(c, "is %T", into)
		}
		*into = JSON{}.Encode(c, from)
		warnIneff(c, "Warn: inefficient json.RawMessage, use enc.Node instead")
		return nil
	case *time.Time:
		var t time.Time
		err := json.Unmarshal(JSON{}.Encode(c, from), &t)
		if err != nil {
			return ctx.NewErrorf(c, "can't unmarshal %#v as time", from)
		}
		*into = t
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

	//log.Debugf(c, "OHA %T => %v", from, v.Type())
	if this.Debugf != nil {
		this.Debugf(c, "normal type: %v, use generic %T.unmarshalInto()", v.Type(), from)
	}
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

func (this Handler) Marshal(c ctx.C, in any) (Node, error) {
	//log.Warnf(c, "OHA: %T %v", in, in)
	switch in := in.(type) {
	case nil:
		return Nil{}, nil
	case Marshaler:
		//log.Warnf(c, "OHA: %T->MarshalNode", in)
		return in.MarshalNode(c)
	case time.Time:
		return Time(in), nil
		//return String(in.Format(time.RFC3339Nano)), nil
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
	case reflect.Func:
		return nil, ctx.NewErrorf(c, "can't marshal %T", in)
	case reflect.Bool:
		return Bool(v.Bool()), nil
	case reflect.Float64, reflect.Float32:
		return Float(v.Float()), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Integer(v.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Integer(v.Uint()), nil
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
		fallthrough //Intentional, since the above check breaks when called on arrays
	case reflect.Array:
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
			switch k := kv.Interface().(type) {
			case string:
				m[k] = n
			default:
				nk, err := this.Marshal(c, k)
				if err != nil {
					return nil, err
				}
				switch nk := nk.(type) {
				case String:
					m[string(nk)] = n
				case Integer:
					m[fmt.Sprint(nk)] = n
				case Float:
					m[fmt.Sprint(nk)] = n
				case Digits:
					m[fmt.Sprint(nk)] = n
				default:
					return nil, ctx.NewErrorf(c, "can't marshal %v as map key", kv.Type())
				}
			}
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
		switch fn := fn.(type) {
		case Nil:
			if !tag.OmitEmpty {
				out = append(out, Pair{tag.Name, tag.JSON, fn})
			}
		case String:
			if !tag.OmitEmpty || fn != "" {
				out = append(out, Pair{tag.Name, tag.JSON, fn})
			}
		case Integer:
			if !tag.OmitEmpty || fn != 0 {
				out = append(out, Pair{tag.Name, tag.JSON, fn})
			}
		case Float:
			if !tag.OmitEmpty || fn != 0.0 {
				out = append(out, Pair{tag.Name, tag.JSON, fn})
			}
		case Digits:
			if !tag.OmitEmpty || fn != "" {
				out = append(out, Pair{tag.Name, tag.JSON, fn})
			}
		default:
			out = append(out, Pair{tag.Name, tag.JSON, fn})
		}
	}
	return out, nil
}
