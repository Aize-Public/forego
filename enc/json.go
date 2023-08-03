package enc

import (
	"bytes"
	"encoding/json"

	"github.com/Aize-Public/forego/ctx"
)

func MustMarshal(c ctx.C, from any) Node {
	n, err := Marshal(c, from)
	if err != nil {
		panic(err)
	}
	return n
}

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
	Indent bool
}

var _ Codec = JSON{}

func (this JSON) Encode(c ctx.C, n Node) []byte {
	if this.Indent {
		j, err := json.MarshalIndent(n, "", "  ")
		if err != nil {
			panic(err)
		}
		return j
	} else {
		j, err := json.Marshal(n)
		if err != nil {
			panic(err)
		}
		return j
	}
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
	dec := json.NewDecoder(bytes.NewBuffer(data))
	dec.UseNumber()
	err := dec.Decode(&obj)
	if err != nil {
		return nil, ctx.NewErrorf(c, "%w", err)
	}
	return fromNative(obj), nil
}
