package ws

import (
	"io"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
)

var _ ws = &TestClient[any]{}

// creates a TestClient which can be used to locally test a client against a websocket implementation
// the websocket will run on its own separate thread
func (this *Handler[State]) TestClient(c ctx.C) *TestClient[State] {
	cli := &TestClient[State]{
		cli2ser: make(chan enc.Node),
		ser2cli: make(chan enc.Node, 100),
	}
	cli.Client = newClient[State](c, func(c ctx.C, n enc.Node) error {
		select {
		case cli.cli2ser <- n:
			return nil
		case <-c.Done():
			return ctx.Cause(c)
		}
	})
	cli.SetFallback(c, func(c ctx.C, n enc.Node) error {
		select {
		case cli.ser2cli <- n:
			return nil
		case <-c.Done():
			return ctx.Cause(c)
		}
	})
	cli.Client.Register(c, this)

	cli.Conn = this.connect(c, cli)
	go func() {
		// run the main loop in a separate go routine
		cli.Conn.loop(c)
	}()
	return cli
}

type TestClient[State any] struct {
	*Client[State]
	Conn    *Conn[State]
	cli2ser chan enc.Node
	ser2cli chan enc.Node
}

func (this *TestClient[State]) Close() error {
	close(this.cli2ser) // closing this channel will provoke a io.EOF on the server side
	close(this.ser2cli)
	return nil
}

func (this *TestClient[State]) Write(c ctx.C, n enc.Node) error {
	return this.Client.onRecv(c, n)
}

func (this *TestClient[State]) Read(c ctx.C) (enc.Node, error) {
	select {

	case <-c.Done():
		return nil, ctx.Cause(c)

	case n, ok := <-this.cli2ser:
		if !ok {
			return nil, io.EOF
		}
		log.Debugf(c, "read %s", n)
		return n, nil
	}
}
