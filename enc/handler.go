package enc

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

// generic object that do the Expand()/Conflate()
type Handler struct {
	Factory map[reflect.Type]func(c ctx.C, n Node) (any, error)

	// called if a field is present in the NodeTree but there is no mapping on the object it's expanded into
	UnhandledFields func(c ctx.C, path Path, n Node) error
}

func (this Handler) Expand(c ctx.C, n Node, into any) error {
	v := reflect.ValueOf(into)
	if v.Kind() != reflect.Pointer {
		return ctx.NewErrorf(c, "expected pointer to expand into, got: %T", into)
	}
	return n.expandInto(c, this, Path{}, v.Elem())
}

func (this Handler) expand(c ctx.C, path Path, from Node, into any) error {
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

func (this Handler) Conflate(c ctx.C, in any) (Node, error) {
	switch in := in.(type) {
	case nil:
		return Nil{}, nil
	case json.Marshaler:
		panic("NIY")
	case Marshaler:
		return in.MarshalTree(c)
	}

	v := reflect.ValueOf(in)
	t := v.Type()
	switch t.Kind() {
	default:
		return fromNative(in), nil
	case reflect.Pointer:
		return this.Conflate(c, v.Elem().Interface())
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
		fn, err := this.Conflate(c, fv.Interface())
		if err != nil {
			return nil, err
		}
		p := Pair{tag.Name, tag.JSON, fn}
		out = append(out, p)
	}
	return out, nil
}
