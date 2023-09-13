package enc

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Aize-Public/forego/ctx"
)

type Time time.Time

var _ Node = Time(time.Time{})

func (this Time) native() any {
	return time.Time(this)
}

func (this Time) GoString() string {
	return fmt.Sprintf("enc.Time{%v}", time.Time(this))
}

func (this Time) String() string {
	return fmt.Sprintf("%v", time.Time(this))
}

func (this Time) unmarshalInto(c ctx.C, handler Handler, into reflect.Value) error {
	v := reflect.ValueOf(this)
	if v.CanConvert(into.Type()) {
		into.Set(v)
		return nil
	}
	return ctx.NewErrorf(c, "can't unmarshal %s %T into %v", handler.path, this, into.Type())
}
