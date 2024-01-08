package ws

import (
	"io"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
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
	if err := enc.Unmarshal(c, n, &f); err != nil {
		return ctx.WrapError(c, err)
	}

	h := this.byChan[f.Channel]
	if h == nil {
		return ctx.NewErrorf(c, "unknown channel %q", f.Channel)
	}
	return h(c, f)
}

// create a local websocket loop and return a client connected to it
func (this *Handler) NewTest(c ctx.C) *TestClient {
	ws := &testWS{
		byChan: map[string]func(ctx.C, Frame) error{},
		inbox:  make(chan Frame, 10),
	}
	conn := &Conn{
		h:  this,
		ws: ws,
	}
	go conn.Loop(c) // nolint

	return &TestClient{
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

// open a channel, and return a call function that will send the data to the given path, and block until
// 1. the function finish with no error
// 2. the function finish with an error
// 3. the callback onData returns an error
func (this *TestClient) Open(c ctx.C, path string, data any, onData func(ctx.C, Frame) error) (call func(ctx.C, string, any) error, err error) {
	chID := uuid.NewString()
	ech := make(chan error, 1)
	this.ws.byChan[chID] = func(c ctx.C, f Frame) error {
		switch f.Type {
		case "return", "error":
			var err error
			if f.Data != nil {
				err = ctx.NewErrorf(c, "remote: %s", f.Data)
			}
			select {
			case ech <- err:
			default:
			}
		}
		err := onData(c, f)
		if err != nil {
			select {
			case ech <- err:
			default:
			}
		}
		return nil
	}
	return func(c ctx.C, path string, data any) error {
			err := this.Send(c, Frame{
				Channel: chID,
				Path:    path,
				Data:    enc.MustMarshal(c, data),
			})
			if err != nil {
				return err
			}
			select {
			case err := <-ech:
				return err
			case <-c.Done():
				return c.Err()
			}
		}, this.Send(c, Frame{
			Type:    "open",
			Channel: chID,
			Path:    path,
			Data:    enc.MustMarshal(c, data),
		})
}

// experimental
// open a channel, and return a ws.C for that channel and a inbox where all the responses will be added to
func (this *TestClient) NewContext(c ctx.C, size int) (C, <-chan Frame) {
	outbox := make(chan Frame, size)
	ch := Channel{
		Conn: this.conn,
		ID:   uuid.NewString(),
	}
	this.ws.byChan[ch.ID] = func(c ctx.C, f Frame) error {
		select {
		case outbox <- f:
			return nil
		case <-c.Done():
			return c.Err()
		}
	}
	return C{
		C:  c,
		ch: &ch,
	}, outbox
}
