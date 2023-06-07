package enc

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Aize-Public/forego/ctx"
)

type List []Node

var _ Node = List{}

func (this List) native() any {
	out := []any{}
	for _, n := range this {
		out = append(out, n.native())
	}
	return out
}

func (this List) GoString() string {
	list := []string{}
	for _, p := range this {
		list = append(list, fmt.Sprintf("%#s", p))
	}
	return "enc.List{" + strings.Join(list, ", ") + "}"
}

func (this List) String() string {
	list := []string{}
	for _, p := range this {
		list = append(list, p.String())
	}
	return "[" + strings.Join(list, ", ") + "]"
}

func (this List) unmarshalInto(c ctx.C, handler Handler, into reflect.Value) error {
	switch into.Kind() {
	case reflect.Interface:
		into.Set(reflect.ValueOf(this.native()))
		return nil

	case reflect.Slice:
		slice := reflect.MakeSlice(into.Type(), len(this), len(this))
		for i := 0; i < len(this); i++ {
			ev := slice.Index(i)
			err := handler.Append(i).unmarshal(c, this[i], ev)
			//err := this[i].unmarshalInto(c, handler, path.Append(i), ev)
			if err != nil {
				return err
			}
		}
		into.Set(slice)
		return nil

	case reflect.Array:
		if len(this) != into.Type().Len() {
			return ctx.NewErrorf(c, "expected %v, got %d elements instead", into.Type(), len(this))
		}
		array := reflect.ArrayOf(len(this), into.Type().Elem())
		instance := reflect.New(array).Elem()
		for i := 0; i < len(this); i++ {
			ev := instance.Index(i)
			err := handler.Append(i).unmarshal(c, this[i], ev)
			if err != nil {
				return err
			}
		}
		into.Set(instance)
		return nil

	case reflect.Invalid:
		return ctx.NewErrorf(c, "can't unmarshal %T into %v at %s", this, into, handler.path.String())

	default:
		return ctx.NewErrorf(c, "can't unmarshal %T into %v at %s", this, into.Type(), handler.path.String())
	}
}
