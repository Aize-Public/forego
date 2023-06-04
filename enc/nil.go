package enc

import (
	"reflect"

	"github.com/Aize-Public/forego/ctx"
)

type Nil struct{}

var _ Node = Nil{}

func (this Nil) native() any { return nil }

func (this Nil) MarshalJSON() ([]byte, error) { return []byte("null"), nil }

func (this Nil) GoString() string {
	return "enc.Nil{}"
}

func (this Nil) String() string {
	return "null"
}

func (this Nil) unmarshalInto(c ctx.C, handler Handler, path Path, into reflect.Value) error {
	into.SetZero()
	return nil
}
