package ws

import (
	"io"
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
	"github.com/google/uuid"
)

type TestClient struct {
	conn *Conn
	ws   *testWS
}

type testWS struct {
	byChan map[string]func(ctx.C, Frame) error
	inbox  chan Frame
	closed bool
}

var _ impl = &testWS{}

func (this *testWS) Close(c ctx.C, reason int) error {
	this.closed = true
	return nil
}

func (this *testWS) Read(c ctx.C) (enc.Node, error) {
	select {
	case f, ok := <-this.inbox:
		if !ok {
			return nil, io.EOF
		}
		log.Debugf(c, "%T.read() %+v", this, f)
		return enc.Marshal(c, f)
	case <-c.Done():
		return nil, c.Err()
	}
}

func (this *testWS) Write(c ctx.C, n enc.Node) error {
	if this.closed {
		return io.ErrClosedPipe
	}
	var f Frame
	enc.MustUnmarshal(c, n, &f)

	h := this.byChan[f.Channel]
	if h == nil {
		return ctx.NewErrorf(c, "unknown channel %q", f.Channel)
	}
	return h(c, f)
}

func (this *Handler) NewTest(t *testing.T) TestClient {
	ws := &testWS{
		byChan: map[string]func(ctx.C, Frame) error{},
		inbox:  make(chan Frame, 10),
	}
	conn := &Conn{
		h:  this,
		ws: ws,
	}
	go conn.Loop(test.Context(t))

	return TestClient{
		conn: conn,
		ws:   ws,
	}
}

func (this *TestClient) Send(c ctx.C, f Frame) error {
	select {
	case this.ws.inbox <- f:
		log.Debugf(c, "%T.send() %+v", this, f)
		return nil
	case <-c.Done():
		return c.Err()
	}
}

func (this *TestClient) Close(c ctx.C) error {
	return this.conn.Close(c, 1000)
}

func (this *TestClient) Open(c ctx.C, path string, data any, onData func(ctx.C, Frame) error) (send func(ctx.C, string, any) error, err error) {
	ch := uuid.NewString()
	this.ws.byChan[ch] = onData
	return func(c ctx.C, path string, data any) error {
			return this.Send(c, Frame{
				Channel: ch,
				Path:    path,
				Data:    enc.MustMarshal(c, data),
			})
		}, this.Send(c, Frame{
			Type:    "open",
			Channel: ch,
			Path:    path,
			Data:    enc.MustMarshal(c, data),
		})
}
