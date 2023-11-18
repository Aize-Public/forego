package ws

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
	"golang.org/x/net/websocket"
)

type Counter struct {
	MinAmt int
	Ct     int
}

func (this *Counter) Init(c C, amt int) error {
	log.Warnf(c, "%p.Init(%v)", this, amt)
	this.Ct = amt
	return nil
}

// name: increment
func (this *Counter) Inc(c C, amt int) error {
	log.Warnf(c, "%p.Inc(%v)", this, amt)
	if amt < this.MinAmt {
		return c.Error(fmt.Sprintf("amt %d < %d", amt, this.MinAmt))
	}
	this.Ct += amt
	return this.Get(c)
}

func (this Counter) Get(c C) error {
	log.Warnf(c, "%p.Get()", &this)
	return c.Reply("ct", this.Ct)
}

func (this Counter) Special(c C, req struct {
	ID  string `json:"id"`
	Val any    `json:"val"`
}) error {
	if req.ID == "" {
		return io.EOF
	}
	return nil
}

func (this Counter) Bye(c C) error {
	return c.Close()
}

func (this *Counter) internal(c C) error { // nolint (meant to be unused)
	return nil
}

func (this *Counter) Foo(other int) {
}

func TestReflect(t *testing.T) {
	var _ Frame
	c := test.Context(t)
	h := Handler{}
	h.MustRegister(c, &Counter{MinAmt: 1})

	send := make(chan chanMsg, 10)
	recv := make(chan chanMsg, 10)
	conn := Conn{
		h: &h,
		ws: &chanImpl{
			Send: send,
			Recv: recv,
		},
	}

	test.NoError(t, conn.onData(c, Frame{
		Channel: "001",
		Path:    "counter",
		Type:    "open",
		Data:    enc.Integer(4),
	}))
	test.NoError(t, conn.onData(c, Frame{
		Channel: "001",
		Path:    "get",
	}))
	test.NoError(t, conn.onData(c, Frame{
		Channel: "001",
		Path:    "inc",
		Data:    enc.Integer(0),
	}))
	test.NoError(t, conn.onData(c, Frame{
		Channel: "001",
		Path:    "inc",
		Data:    enc.Integer(3),
	}))
	test.NoError(t, conn.onData(c, Frame{
		Channel: "001",
		Path:    "get",
	}))
	test.NoError(t, conn.onData(c, Frame{
		Channel: "001",
		Path:    "bye",
	}))
	test.NoError(t, conn.Close(c, 1000))

	time.Sleep(time.Millisecond)
	for msg := range send {
		switch msg.Type {
		case websocket.TextFrame:
			var f Frame
			enc.MustUnmarshal(c, msg.Data, &f)
			if f.Type == "error" {
				test.ContainsJSON(t, f.Data, "amt")
			} else {
				test.OK(t, "recv: %+v", f.Data)
			}
		case websocket.CloseFrame:
			t.Logf("CLOSE")
		default:
			test.Fail(t, "unexpected %v", msg)
		}
	}
	t.Logf("EXIT")
}
