package enc

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

// unordered map
type Map map[string]Node

var _ Node = Map{}

func (this Map) native() any {
	out := map[string]any{}
	for k, n := range this {
		out[k] = n.native()
	}
	return out
}

func (this Map) String() string {
	p := []string{}
	for k, v := range this {
		p = append(p, fmt.Sprintf("%q:%#s", k, v))
	}
	return "{" + strings.Join(p, ", ") + "}"
}

func (this Map) GoString() string {
	p := []string{}
	for k, v := range this {
		p = append(p, fmt.Sprintf("%q:%#s", k, v))
	}
	return "enc.Map{" + strings.Join(p, ", ") + "}"
}

func (this Map) unmarshalInto(c ctx.C, handler Handler, path Path, into reflect.Value) error {
	switch into.Kind() {
	case reflect.Map:
		log.Debugf(c, "%T.expandInto() => map", this)
		t := into.Type()
		intokt := t.Key()
		intovt := t.Elem()
		mv := reflect.MakeMap(t)
		for k, v := range this {
			kv := reflect.New(intokt)
			err := handler.unmarshal(c, path.Append("#key"), String(k), kv.Interface())
			if err != nil {
				return err
			}
			vv := reflect.New(intovt)
			err = handler.unmarshal(c, path.Append("#val"), v, vv.Interface())
			if err != nil {
				return err
			}
			mv.SetMapIndex(kv.Elem(), vv.Elem())
		}
		into.Set(mv)
		return nil

	case reflect.Struct:
		seen := map[string]bool{}
		vt := into.Type()
		for i := 0; i < into.NumField(); i++ {
			fv := into.Field(i)
			ft := vt.Field(i)
			if ft.IsExported() {
				tag := parseTag(ft)
				seen[tag.JSON] = true
				v, ok := this[tag.JSON]
				if ok {
					err := handler.unmarshal(c, path.Append(tag.Name), v, fv.Addr().Interface())
					if err != nil {
						return err
					}
				}
			}
		}
		if handler.UnhandledFields != nil {
			for k, v := range this {
				if !seen[k] {
					err := handler.UnhandledFields(c, append(path, k), v)
					if err != nil {
						return ctx.NewError(c, err)
					}
				}
			}
		}
		return nil
	case reflect.Interface:
		into.Set(reflect.ValueOf(this.native()))
		return nil

	default:
		return ctx.NewErrorf(c, "can't expand %T into %v", this, into.Type())
	}
}
