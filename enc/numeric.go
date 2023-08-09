package enc

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/Aize-Public/forego/ctx"
)

type Numeric interface {
	Int64() (int64, error)
	Float64() (float64, error)
	String() string
}
type numeric interface {
	json.Marshaler
	Numeric
	Node
}

type Integer int

var _ Node = Integer(0)
var _ Numeric = Integer(0)

func (this Integer) Int64() (int64, error)        { return int64(this), nil }
func (this Integer) Float64() (float64, error)    { return float64(this), nil }
func (this Integer) String() string               { return strconv.FormatInt(int64(this), 10) }
func (this Integer) GoString() string             { return fmt.Sprintf("enc.Int{%v}", int64(this)) }
func (this Integer) native() any                  { return int64(this) }
func (this Integer) MarshalJSON() ([]byte, error) { return json.Marshal(int64(this)) }
func (this Integer) unmarshalInto(c ctx.C, handler Handler, into reflect.Value) error {
	return unmarshalNumericInto(this, c, handler, into)
}

type Float float64

var _ Node = Float(0)
var _ Numeric = Float(0)

func (this Float) Int64() (int64, error)        { return int64(this), nil }
func (this Float) Float64() (float64, error)    { return float64(this), nil }
func (this Float) String() string               { return strconv.FormatFloat(float64(this), 'g', -1, 64) }
func (this Float) GoString() string             { return fmt.Sprintf("enc.Float{%v}", float64(this)) }
func (this Float) native() any                  { return float64(this) }
func (this Float) MarshalJSON() ([]byte, error) { return json.Marshal(float64(this)) }
func (this Float) unmarshalInto(c ctx.C, handler Handler, into reflect.Value) error {
	return unmarshalNumericInto(this, c, handler, into)
}

type Digits string

var _ numeric = Digits("0")

func (this Digits) Int64() (int64, error) {
	i, err := strconv.ParseInt(string(this), 10, 64)
	return i, err
}
func (this Digits) Float64() (float64, error) { return strconv.ParseFloat(string(this), 64) }
func (this Digits) String() string            { return string(this) }
func (this Digits) GoString() string          { return fmt.Sprintf("enc.Num{%q}", string(this)) }

func (this Digits) native() any {
	f, err := this.Float64() // builtin json convert to float64 when unmarshalling into any, we should do the same
	if err == nil {
		return f
	}
	return this.String()
}
func (this Digits) MarshalJSON() ([]byte, error) { return []byte(this), nil }
func (this Digits) unmarshalInto(c ctx.C, handler Handler, into reflect.Value) error {
	return unmarshalNumericInto(this, c, handler, into)
}

func unmarshalNumericInto(this numeric, c ctx.C, handler Handler, into reflect.Value) error {
	//log.Debugf(c, "OHA %#v => %#v", this, into)
	switch into.Kind() {
	case reflect.Float64, reflect.Float32:
		f, err := this.Float64()
		if err != nil {
			return err
		}
		into.SetFloat(f)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		i, err := this.Int64()
		if err != nil {
			return err
		}
		into.SetInt(int64(i))
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		i, err := this.Int64() // TODO FIXME need Uint
		if err != nil {
			return err
		}
		into.SetUint(uint64(i))
	case reflect.Interface:
		v := reflect.ValueOf(this.native())
		into.Set(v)
	default:
		return ctx.NewErrorf(c, "can't unmarshal %s %T into %v", handler.path, this, into.Type())
	}
	return nil
}
