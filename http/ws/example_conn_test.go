package ws_test

import (
	"testing"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/http"
	"github.com/Aize-Public/forego/http/ws"
	"github.com/Aize-Public/forego/test"
)

func TestExampleConn(t *testing.T) {
	c := test.Context(t)

	s := http.NewServer(c)
	s.Handle("/ws", ws.Handler{
		OnConnect: func(c ctx.C, conn *ws.Conn) (ctx.C, error) {
			return c, conn.Send(c, "hello there!")
		},
		OnData: func(c ctx.C, conn *ws.Conn, n enc.Node) error {
			return conn.Send(c, n) // echo
		},
	}.Server())
}
