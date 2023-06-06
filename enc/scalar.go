package enc

import (
	"fmt"
	"reflect"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

type String string

var _ Node = String("")

func (this String) native() any {
	return string(this)
}

func (this String) GoString() string {
	return fmt.Sprintf("enc.String{%q}", string(this))
}

func (this String) String() string {
	return fmt.Sprintf("%q", string(this))
}

func (this String) unmarshalInto(c ctx.C, handler Handler, into reflect.Value) error {
	//log.Debugf(c, "%v.unmarshalInto(%#v)", this, into)
	switch into.Kind() {
	case reflect.String:
		into.SetString(string(this))
	case reflect.Interface:
		v := reflect.ValueOf(this.native())
		into.Set(v)
	default:
		return ctx.NewErrorf(c, "can't unmarshal %s %T into %v", handler.path, this, into.Type())
	}
	return nil
}

// NOTE(oha) do we need to split into Integers and Floats?
type Number float64

var _ Node = Number(1.0)

func (this Number) native() any {
	return float64(this)
}

func (this Number) GoString() string {
	return fmt.Sprintf("enc.Number{%v}", float64(this))
}

func (this Number) String() string {
	return fmt.Sprintf("%v", float64(this))
}

func (this Number) unmarshalInto(c ctx.C, handler Handler, into reflect.Value) error {
	switch into.Kind() {
	case reflect.Float64, reflect.Float32:
		into.SetFloat(float64(this))
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		into.SetInt(int64(this))
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		into.SetUint(uint64(this))
	case reflect.Interface:
		v := reflect.ValueOf(this.native())
		into.Set(v)
	default:
		return ctx.NewErrorf(c, "can't unmarshal %T into %v", this, into.Type())
	}
	return nil
}

type Bool bool

var _ Node = Bool(true)

func (this Bool) native() any {
	return bool(this)
}

func (this Bool) GoString() string {
	if this {
		return "enc.Bool{true}"
	} else {
		return "enc.Bool{false}"
	}
}

func (this Bool) String() string {
	if this {
		return "true"
	} else {
		return "false"
	}
}

func (this Bool) unmarshalInto(c ctx.C, handler Handler, into reflect.Value) error {
	log.Debugf(c, "%v.unmarshalInto(%#v)", this, into)
	switch into.Kind() {
	case reflect.Bool:
		into.SetBool(bool(this))
	case reflect.Interface:
		into.Set(reflect.ValueOf(bool(this)))
	default:
		return ctx.NewErrorf(c, "can't unmarshal %T into %v", this, into.Type())
	}
	return nil
}
