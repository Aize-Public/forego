package ws

import (
	"io"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"golang.org/x/net/websocket"
)

type chanMsg struct {
	Type int
	Data enc.Node
}

type chanImpl struct {
	Send chan<- chanMsg
	Recv <-chan chanMsg
}

func (this chanImpl) Write(c ctx.C, n enc.Node) error {
	log.Debugf(c, "%p: sending %+v", this.Send, n)
	select {
	case <-c.Done():
		return c.Err()
	case this.Send <- chanMsg{
		Type: websocket.TextFrame,
		Data: n,
	}:
		return nil
	}
}

func (this chanImpl) Read(c ctx.C) (enc.Node, error) {
	select {
	case <-c.Done():
		return nil, c.Err()
	case msg, ok := <-this.Recv:
		if !ok {
			return nil, ctx.NewErrorf(c, "ws.exit<closed>")
		}
		switch msg.Type {
		case websocket.CloseFrame:
			switch msg.Data {
			case enc.Integer(EXIT):
				return nil, io.EOF
			default:
				return nil, ctx.NewErrorf(c, "ws.exit<%v>", msg.Data)
			}
		case websocket.TextFrame:
			return msg.Data, nil
		default:
			return nil, ctx.NewErrorf(c, "unexpected type: %v", msg)
		}
	}
}

func (this chanImpl) Close(c ctx.C, reason int) error {
	time.Sleep(time.Millisecond)
	log.Debugf(c, "%p: closing %+v", this.Send, reason)
	select {
	case <-c.Done():
		return c.Err()
	case this.Send <- chanMsg{
		Type: websocket.CloseFrame,
		Data: enc.Integer(reason),
	}:
		close(this.Send)
		return nil
	}
}
