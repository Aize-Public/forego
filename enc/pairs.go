package enc

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/Aize-Public/forego/ctx"
)

// ordered pairs, used mostly internally when Unmarshalling structs, to presever the order of the fields
// can be used anywhere else where the order matters
// Note(oha): currently not used while decoding, to keep the decode stack simple and easy, might be changed in the future?
type Pairs []Pair

var _ Node = Pairs{}

func (this Pairs) native() any {
	out := map[string]any{}
	for _, p := range this {
		out[p.Name] = p.Value.native()
	}
	return out
}

func (this Pairs) String() string {
	list := []string{}
	for _, p := range this {
		list = append(list, fmt.Sprintf("%q{json:%q}:%#s", p.Name, p.JSON, p.Value))
	}
	return "enc.Pairs{" + strings.Join(list, ", ") + "}"
}

func (this Pairs) GoString() string {
	list := []string{}
	for _, p := range this {
		list = append(list, fmt.Sprintf("%q:%#s", p.Name, p.Value))
	}
	return "enc.Pairs{" + strings.Join(list, ", ") + "}"
}

func (this Pairs) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")
	for i, p := range this {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.Write(mustJSON(p.JSON))
		buf.WriteString(": ")
		buf.Write(mustJSON(p.Value))
	}
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (this Pairs) Find(name string) Node {
	for _, p := range this {
		if p.Name == name {
			return p.Value
		}
	}
	return nil
}

func (this Pairs) unmarshalInto(c ctx.C, handler Handler, path Path, into reflect.Value) error {
	panic("NIY")
}

type Pair struct {
	Name  string
	JSON  string
	Value Node
}

func (this Pair) jsonName() string {
	if this.JSON != "" {
		return this.JSON
	}
	return this.Name
}
