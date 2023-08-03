package enc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Aize-Public/forego/ctx"
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

func (this Map) MarshalJSON() ([]byte, error) {
	if this == nil {
		return []byte(`{}`), nil
	}
	return json.Marshal(map[string]Node(this))
}

func (this Map) String() string {
	p := []string{}
	for k, v := range this {
		p = append(p, fmt.Sprintf("%q:%s", k, v))
	}
	return "{" + strings.Join(p, ", ") + "}"
}

func (this Map) GoString() string {
	p := []string{}
	for k, v := range this {
		p = append(p, fmt.Sprintf("%q:%s", k, v))
	}
	return "enc.Map{" + strings.Join(p, ", ") + "}"
}

func (this Map) unmarshalInto(c ctx.C, handler Handler, into reflect.Value) error {
	switch into.Kind() {
	case reflect.Map:
		t := into.Type()
		intokt := t.Key()
		intovt := t.Elem()
		mv := reflect.MakeMap(t)
		for k, n := range this {
			kv := reflect.New(intokt).Elem()
			switch intokt.Kind() {
			case reflect.String:
				kv.SetString(k)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				i, err := strconv.ParseInt(k, 10, 64)
				if err != nil {
					return ctx.NewErrorf(c, "can't convert %q to string at %v", k, handler.path.String()+"#key")
				}
				kv.SetInt(i)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				i, err := strconv.ParseUint(k, 10, 64)
				if err != nil {
					return ctx.NewErrorf(c, "can't convert %q to string at %v", k, handler.path.String()+"#key")
				}
				kv.SetUint(i)
			default:
				err := handler.Append("#key").unmarshal(c, String(k), kv)
				if err != nil {
					return ctx.NewErrorf(c, "can't convert %q to %v", k, intokt)
				}
			}
			vv := reflect.New(intovt).Elem()
			err := handler.Append("#val").unmarshal(c, n, vv)
			if err != nil {
				return err
			}
			mv.SetMapIndex(kv, vv)
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
					err := handler.Append(tag.Name).unmarshal(c, v, fv)
					if err != nil {
						return err
					}
				}
			}
		}
		if handler.UnhandledFields != nil {
			for k, v := range this {
				if !seen[k] {
					err := handler.UnhandledFields(c, append(handler.path, k), v)
					if err != nil {
						return ctx.WrapError(c, err)
					}
				}
			}
		}
		return nil
	case reflect.Interface:
		if handler.Debugf != nil {
			handler.Debugf(c, "assign %v to %v", into.Type(), this)
		}
		into.Set(reflect.ValueOf(this.native()))
		return nil

	default:
		return ctx.NewErrorf(c, "can't expand %s %T into %v", handler.path, this, into.Type())
	}
}
