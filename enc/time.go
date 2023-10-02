package enc

import (
	"encoding/json"
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

// Use this object if you want to get `1s` from time.Second
// the built in json library encode time.Second as 1000000000
type Duration time.Duration

var _ Node = Duration(time.Second)

func (this Duration) native() any {
	return time.Duration(this)
}

func (this Duration) GoString() string {
	return fmt.Sprintf("enc.Duration{%v}", time.Duration(this))
}

func (this Duration) String() string {
	return fmt.Sprintf("%v", time.Duration(this))
}

func (this Duration) unmarshalInto(c ctx.C, handler Handler, into reflect.Value) error {
	v := reflect.ValueOf(this)
	if v.CanConvert(into.Type()) {
		into.Set(v)
		return nil
	}
	return ctx.NewErrorf(c, "can't unmarshal %s %T into %v", handler.path, this, into.Type())
}

func (this Duration) MarshalJSON() ([]byte, error) {
	if this == 0 {
		return json.Marshal(0)
	}
	return json.Marshal(this.String())
}

func (this *Duration) UnmarshalJSON(in []byte) error {
	var s string
	err := json.Unmarshal(in, &s)
	if err != nil {
		return err
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*this = Duration(d)
	return nil
}
