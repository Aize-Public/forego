package enc

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Aize-Public/forego/ctx"
)

type Node interface {
	native() any
	unmarshalInto(c ctx.C, handler Handler, into reflect.Value) error
	String() string
	GoString() string
}

type Codec interface {
	Encode(c ctx.C, n Node) []byte
	Decode(c ctx.C, data []byte) (Node, error)
}

// objects which implements this can override how their data is unmarshaled (Expand)
type Unmarshaler interface {
	UnmarshalNode(ctx.C, Node) error
}

// obejcts which implements this can override how they are marshaled (Conflate)
type Marshaler interface {
	MarshalNode(ctx.C) (Node, error)
}

type path []any

func (this path) String() string {
	out := ""
	for _, v := range this {
		switch v := v.(type) {
		case string:
			out += "." + v
		case int:
			out += fmt.Sprintf("[%d]", v)
		default:
			out += fmt.Sprintf("{%v}", v)
		}
	}
	if out == "" {
		return "ROOT"
	}
	return strings.TrimPrefix(out, ".")
}

func (this path) Parent() path {
	if len(this) > 0 {
		return this[0 : len(this)-2]
	}
	return this
}

func fromNative(in any) Node {
	switch in := in.(type) {
	case nil:
		return Nil{}
	case map[string]any:

		out := Map{}
		for k, v := range in {
			out[k] = fromNative(v)
		}
		return out
	case []any:
		out := List{}
		for _, v := range in {
			out = append(out, fromNative(v))
		}
		return out

	case string:
		return String(in)

	case float64:
		return Number(in)
	case float32:
		return Number(in)

	case int:
		return Number(in)
	case int8:
		return Number(in)
	case int16:
		return Number(in)
	case int32:
		return Number(in)
	case int64:
		return Number(in)
	case uint:
		return Number(in)
	case uint8:
		return Number(in)
	case uint16:
		return Number(in)
	case uint32:
		return Number(in)
	case uint64:
		return Number(in)

	case bool:
		return Bool(in)
	default:
		panic(fmt.Sprintf("unexpected native type %T: %+v", in, in))
	}
}
