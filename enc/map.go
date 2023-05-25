package enc

import (
	"bytes"
	"fmt"
	"reflect"
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

func (this Map) expandInto(c ctx.C, codec Codec, path Path, into reflect.Value) error {
	switch into.Kind() {
	case reflect.Map:
		t := into.Type()
		intokt := t.Key()
		intovt := t.Elem()
		mv := reflect.MakeMap(t)
		for k, v := range this {
			kv := reflect.New(intokt)
			err := codec.expand(c, path.Append("#key"), String(k), kv.Interface())
			if err != nil {
				return err
			}
			vv := reflect.New(intovt)
			err = codec.expand(c, path.Append("#val"), v, vv.Interface())
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
				seen[tag.Name] = true
				v, ok := this[tag.Name]
				if ok {
					err := codec.expand(c, path.Append(tag.Name), v, fv.Addr().Interface())
					if err != nil {
						return err
					}
				}
			}
		}
		if codec.UnhandledFields != nil {
			for k, v := range this {
				if !seen[k] {
					err := codec.UnhandledFields(c, append(path, k), v)
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

// ordered pairs, used mostly internally when Unmarshalling structs, to presever the order of the fields
// can be used anywhere else where the order matters
// Note(oha): currently not used while decoding, to keep the decode stack simple and easy, might be changed in the future?
type Pairs []Pair

var _ Node = Pairs{}

func (this Pairs) native() any {
	out := map[string]any{}
	for _, p := range this {
		out[p.Key] = p.Value.native()
	}
	return out
}

func (this Pairs) String() string {
	list := []string{}
	for _, p := range this {
		list = append(list, fmt.Sprintf("%q:%#s", p.Key, p.Value))
	}
	return "enc.Pairs{" + strings.Join(list, ", ") + "}"
}

func (this Pairs) GoString() string {
	list := []string{}
	for _, p := range this {
		list = append(list, fmt.Sprintf("%q:%#s", p.Key, p.Value))
	}
	return "enc.Pairs{" + strings.Join(list, ", ") + "}"
}

type Pair struct {
	Key   string
	Value Node
}

func (this Pairs) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, p := range this {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.Write(mustJSON(p.Key))
		buf.WriteString(": ")
		buf.Write(mustJSON(p.Value))
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (this Pairs) Find(key string) Node {
	for _, p := range this {
		if p.Key == key {
			return p.Value
		}
	}
	return nil
}

func (this Pairs) expandInto(c ctx.C, codec Codec, path Path, into reflect.Value) error {
	panic("NIY")
}
