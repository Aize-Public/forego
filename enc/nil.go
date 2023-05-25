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
	return "<nil>"
}

func (this Nil) expandInto(c ctx.C, codec Codec, path Path, into reflect.Value) error {
	into.SetZero()
	return nil
}
