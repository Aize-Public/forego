package enc

import (
	"fmt"
	"reflect"

	"github.com/Aize-Public/forego/ctx"
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

func (this String) expandInto(c ctx.C, codec Codec, path Path, into reflect.Value) error {
	//log.Debugf(c, "%v.expandInto(%#v)", this, into)
	switch into.Kind() {
	case reflect.String:
		into.SetString(string(this))
	case reflect.Interface:
		v := reflect.ValueOf(this.native())
		into.Set(v)
	default:
		return ctx.NewErrorf(c, "can't expand %T into %v", this, into.Type())
	}
	return nil
}

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

func (this Number) expandInto(c ctx.C, codec Codec, path Path, into reflect.Value) error {
	//log.Debugf(c, "%v.expandInto(%#v)", this, into)
	switch into.Kind() {
	case reflect.Float64:
		into.SetFloat(float64(this))
	case reflect.Interface:
		v := reflect.ValueOf(this.native())
		into.Set(v)
	default:
		return ctx.NewErrorf(c, "can't expand %T into %v", this, into.Type())
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

func (this Bool) expandInto(c ctx.C, codec Codec, path Path, into reflect.Value) error {
	//log.Debugf(c, "%v.expandInto(%#v)", this, into)
	switch into.Kind() {
	case reflect.Bool:
		into.SetBool(bool(this))
	case reflect.Interface:
		v := reflect.ValueOf(bool(this))
		into.Set(v)
	default:
		return ctx.NewErrorf(c, "can't expand %T into %v", this, into.Type())
	}
	return nil
}
