package enc

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/Aize-Public/forego/ctx"
)

func MustMap(n Node) Map {
	switch n := n.(type) {
	case Map:
		return n
	case Pairs:
		return n.AsMap()
	default:
		panic(fmt.Sprintf("not a map: %T", n))
	}
}

func AsMap(c ctx.C, n Node) (Map, error) {
	switch n := n.(type) {
	case Map:
		return n, nil
	case Pairs:
		return n.AsMap(), nil
	default:
		return nil, ctx.NewErrorf(c, "not a map: %T", n)
	}
}

type Pipe struct {
	remoteClose chan struct{}
	Send        chan<- Node
	Recv        <-chan Node
}

var _ ReadWriteCloser = Pipe{}

func NewPipe(buf int) (Pipe, Pipe) {
	send := make(chan Node, buf)
	recv := make(chan Node, buf)
	return Pipe{
			make(chan struct{}),
			send, recv,
		},
		Pipe{
			make(chan struct{}),
			recv, send,
		}
}

func (this Pipe) Read(c ctx.C) (Node, error) {
	select {
	case n, ok := <-this.Recv:
		if !ok {
			return nil, io.EOF
		}
		return n, nil
	case <-c.Done():
		return nil, c.Err()
	}
}

func (this Pipe) Write(c ctx.C, n Node) (err error) {
	defer func() {
		if recover() != nil { // ugly, but simplier this way, allow to write on closed and get an error not a panic
			err = io.ErrClosedPipe
		}
	}()
	select {
	case this.Send <- n:
		return nil
	case <-c.Done():
		return c.Err()
	}
}

func (this Pipe) Close(c ctx.C) error {
	defer func() {
		recover() // ugly, but simplier this way, allowing multiple closes with no side effects
	}()
	close(this.Send)
	return nil
}

type Tag struct {
	Name      string
	JSON      string
	OmitEmpty bool
	Skip      bool
}

// TODO(oha): we need to parse `enc`, `json` and eventually `yaml` and make sure they agree
func parseTag(tag reflect.StructField) (out Tag) {
	out.Name = tag.Name

	json := tag.Tag.Get("json")
	if json == "-" {
		out.Skip = true
		return
	}
	json, extra, _ := strings.Cut(json, ",")
	out.JSON = json
	if out.JSON == "" {
		out.JSON = out.Name
	}
	switch extra {
	case "omitempty":
		out.OmitEmpty = true
	case "":
	default:
		panic(fmt.Sprintf("invalid tag: %v", tag))
	}
	return
}
