package ws_test

import (
	"testing"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/http/ws"
	"github.com/Aize-Public/forego/test"
)

func TestWS(t *testing.T) {
	c := test.Context(t)
	h, err := ws.New[chan bool](c, &Sleep{})
	h.OnConnect = func(c ctx.C, conn *ws.Conn[chan bool]) {
		conn.State = make(chan bool)
	}
	test.NoError(t, err)
	cli := h.TestClient(c)
	defer cli.Close()

	recv := make(chan ws.Op[chan bool])
	r, err := cli.Request(c, &Sleep{}, func(c ctx.C, op ws.Op[chan bool]) error {
		recv <- op
		return nil
	})
	test.NoError(t, err)
	defer r.Close()

	test.Empty(t, recv)

	cli.Conn.State <- true // unblock
	select {
	case op := <-recv:
		test.NotEmpty(t, op)
	case <-time.After(time.Second):
		t.Fatalf("timeout waiting for reply")
	}
}

type Sleep struct {
}

func (this *Sleep) Do(c ctx.C, state chan bool) error {
	<-state // wait for it
	return nil
}

/*
func TestMultiClient(t *testing.T) {
	c := test.Context(t)
	h, err := ws.New[*int](c, &Inc{}, &Get{})
	test.NoError(t, err)

	h.TestClient(c)
}
*/

type Inc struct {
}

func (this *Inc) Do(c ctx.C, state *int) error {
	*state++
	return nil
}

type Get struct {
	Value int `api:"out" json:"value"`
}

func (this *Get) Do(c ctx.C, state *int) error {
	this.Value = *state
	return nil
}
