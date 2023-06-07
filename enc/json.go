package enc

import (
	"encoding/json"

	"github.com/Aize-Public/forego/ctx"
)

func MustMarshalJSON(c ctx.C, from any) []byte {
	n, err := Marshal(c, from)
	if err != nil {
		panic(err)
	}
	return JSON{}.Encode(c, n)
}

func MarshalJSON(c ctx.C, from any) ([]byte, error) {
	n, err := Marshal(c, from)
	if err != nil {
		return nil, err
	}
	return JSON{}.Encode(c, n), nil
}

func UnmarshalJSON(c ctx.C, j []byte, into any) error {
	n, err := JSON{}.Decode(c, j)
	if err != nil {
		return err
	}
	return Unmarshal(c, n, into)
}

type JSON struct {
}

var _ Codec = JSON{}

func (this JSON) Encode(c ctx.C, n Node) []byte {
	j, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	return j
}

func mustJSON(in any) []byte {
	j, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return j
}

func (this JSON) Decode(c ctx.C, data []byte) (Node, error) {
	var obj any
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return nil, ctx.NewError(c, err)
	}
	return fromNative(obj), nil
}
