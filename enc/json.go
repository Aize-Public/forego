package enc

import (
	"encoding/json"

	"github.com/Aize-Public/forego/ctx"
)

type JSON struct {
}

var _ Encoder = JSON{}
var _ Decoder = JSON{}

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
