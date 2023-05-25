package enc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

type Node interface {
	native() any
	expandInto(c ctx.C, codec Codec, path Path, into reflect.Value) error
}

type Encoder interface {
	Encode(c ctx.C, n Node) []byte
}

type Decoder interface {
	Decode(c ctx.C, data []byte) (Node, error)
}

// generic object that do the Expand()/Conflate()
type Codec struct {
	Factory map[reflect.Type]func(c ctx.C, n Node) (any, error)

	// called if a field is present in the NodeTree but there is no mapping on the object it's expanded into
	UnhandledFields func(c ctx.C, path Path, n Node) error
}

func (this Codec) Expand(c ctx.C, n Node, into any) error {
	v := reflect.ValueOf(into)
	if v.Kind() != reflect.Pointer {
		return ctx.NewErrorf(c, "expected pointer to expand into, got: %T", into)
	}
	return n.expandInto(c, this, Path{}, v.Elem())
}

type Unmarshaler interface {
	UnmarshalTree(ctx.C, Node) error
}

type Marshaler interface {
	MarshalTree(ctx.C) (Node, error)
}

type Path []any

func (this Path) String() string {
	out := ""
	for _, v := range this {
		switch v := v.(type) {
		case string:
			out += "." + v
		case int:
			out += fmt.Sprintf("[%d]", v)
		}
	}
	if out == "" {
		return "ROOT"
	}
	return strings.TrimPrefix(out, ".")
}

func (this Path) Append(d any) Path {
	return append(this, d)
}

func (this Path) Parent() Path {
	if len(this) > 0 {
		return this[0 : len(this)-2]
	}
	return this
}

func (this Codec) expand(c ctx.C, path Path, from Node, into any) error {
	switch into := into.(type) {
	case *time.Time: // override from default go
		switch from := from.(type) {
		case String:
			var err error
			*into, err = time.Parse(time.RFC3339Nano, from.String())
			return err
		default:
			return ctx.NewErrorf(c, "can't expend time.Time from %T (%v)", from, from)
		}
	case Unmarshaler:
		return into.UnmarshalTree(c, from)
	case *json.RawMessage:
		*into = JSON{}.Encode(c, from)
		log.Infof(c, "Warn: inefficient json.RawMessage, use enc.Tree instead")
		return nil
	case json.Unmarshaler:
		log.Infof(c, "Warn: inefficient %T.UnmarshalJSON(): implement %T instead", into, (Unmarshaler)(nil))
		j, _ := json.Marshal(from) // we must go back to the json
		return into.UnmarshalJSON(j)
	case *Node:
		log.Warnf(c, "WTF?")
		*into = from
		return nil
	}
	v := reflect.ValueOf(into)

	if v.Kind() != reflect.Pointer {
		return ctx.NewErrorf(c, "expected pointer, got %T", into)
	}

	vv := v.Elem()
	if from == nil {
		vv.SetZero()
		return nil
	}

	if this.Factory != nil {
		f := this.Factory[v.Type().Elem()] // Note(oha): we get a pointer to the interface to assign to
		if f != nil {
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

	if vv.Kind() == reflect.Pointer {
		// if we expand into a pointer, we create the zero value and expand into that instead
		// this allow to expand into *bool or *string
		e := reflect.New(vv.Type().Elem())
		err := from.expandInto(c, this, path, e.Elem())
		if err != nil {
			return err
		}
		vv.Set(e)
		return nil
	}

	return from.expandInto(c, this, path, vv)
}

func (this Codec) Conflate(c ctx.C, n Node) (Node, error) {
	panic("NIY")
}

func toJson(in Node) []byte {
	j, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return j
}

func fromNative(in any) Node {
	switch in := in.(type) {
	case nil:
		return Nil{}
	case map[string]any:
		out := Map{}
		for k, v := range in {
			out[k] = fromNative(v)
		}
		return out
	case []any:
		out := List{}
		for _, v := range in {
			out = append(out, fromNative(v))
		}
		return out
	case string:
		return String(in)
	case float64:
		return Number(in)
	case bool:
		return Bool(in)
	default:
		panic(fmt.Sprintf("unexpected native type %T: %+v", in, in))
	}
}
