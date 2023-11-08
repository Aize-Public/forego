package ws

import (
	"io"
	"testing"
	"time"

	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

type Counter struct {
	Ct int
}

func (this *Counter) Init(c C, amt int) error {
	log.Warnf(c, "%p.Init(%v)", this, amt)
	this.Ct = amt
	return nil
}

// name: increment
func (this *Counter) Inc(c C, amt int) error {
	log.Warnf(c, "%p.Inc(%v)", this, amt)
	this.Ct += amt
	return this.Get(c)
}

func (this Counter) Get(c C) error {
	log.Warnf(c, "%p.Get()", &this)
	c.Reply("ct", this.Ct)
	return nil
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

func (this *Counter) internal(c C) error {
	return nil
}

func (this *Counter) Foo(other int) {
}

func TestReflect(t *testing.T) {
	var _ Frame
	c := test.Context(t)
	h := Handler{}
	h.Register(c, &Counter{})

	send := make(chan chanMsg, 10)
	recv := make(chan chanMsg, 10)
	conn := Conn{
		h: &h,
		ws: &chanImpl{
			Send: send,
			Recv: recv,
		},
	}

	conn.onData(c, Frame{
		Channel: "001",
		Path:    "counter",
		Data:    enc.Integer(4),
	})
	conn.onData(c, Frame{
		Channel: "001",
		Path:    "get",
	})
	conn.onData(c, Frame{
		Channel: "001",
		Path:    "inc",
		Data:    enc.Integer(3),
	})
	conn.onData(c, Frame{
		Channel: "001",
		Path:    "get",
	})
	conn.Close(c, 1000)

	time.Sleep(time.Millisecond)
	for msg := range send {
		t.Logf("resp: %+v", msg)
	}
	t.Logf("EXIT")
}
