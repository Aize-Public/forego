package ws

import (
	"fmt"
	"net/http"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/shutdown"
	"github.com/Aize-Public/forego/utils/sync"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

// a server side handler for wsrpc requests
type RpcHandler struct {
	byPath sync.Map[string, func(ctx.C, *RpcConn, Frame) error]
}

// a server side handler for generic connections
type Handler struct {
	OnConnect func(c ctx.C, conn *Conn) (ctx.C, error)
	OnData    func(c ctx.C, conn *Conn, data enc.Node) error
}

func (this Handler) Server() websocket.Server {
	if this.OnData == nil {
		panic("Handler.OnData is nil")
	}
	x := websocket.Server{
		Handler: websocket.Handler(func(conn *websocket.Conn) {
			c := conn.Request().Context()
			c, cf := ctx.Span(c, "ws")
			defer cf(nil)

			defer shutdown.Hold().Release()

			//defer metrics.WS{Path: path}.Start().End(c)
			ws := &Conn{
				sid: uuid.NewString(),
				ws: &wsImpl{
					conn: conn,
				},
			}
			ws.onData = func(c ctx.C, data enc.Node) error {
				return this.OnData(c, ws, data)
			}
			c = ctx.WithTag(c, "ws.sid", ws.sid)
			if this.OnConnect != nil {
				var err error
				c, err = this.OnConnect(c, ws)
				if err != nil {
					log.Errorf(c, "can't connect ws: %v", err)
					return
				}
			}
			defer ws.Close(c, 1000)
			err := ws.Loop(c)
			if err != nil {
				log.Warnf(c, "loop: %v", err)
			}
		}),
		Handshake: func(config *websocket.Config, req *http.Request) (err error) {
			config.Origin, err = websocket.Origin(config, req)
			if err == nil && config.Origin == nil {
				return fmt.Errorf("null origin")
			}
			return err
		},
	}
	return x
}

// return a websocket.Server which can be used as an http.Handler
// Note: it sets a default Handshake handler which accept any requests,
// you might need to change it if you need to control the `Origin` header.
func (this *RpcHandler) Server() websocket.Server {
	x := websocket.Server{
		Handler: websocket.Handler(func(conn *websocket.Conn) {
			c := conn.Request().Context()
			c, cf := ctx.Span(c, "ws")
			defer cf(nil)

			defer shutdown.Hold().Release()

			//defer metrics.WS{Path: path}.Start().End(c)
			ws := this.newRpcConn(&wsImpl{conn: conn})
			defer ws.Close(c, 1000)
			err := ws.Loop(c)
			if err != nil {
				log.Warnf(c, "loop: %v", err)
			}
		}),
		Handshake: func(config *websocket.Config, req *http.Request) (err error) {
			config.Origin, err = websocket.Origin(config, req)
			if err == nil && config.Origin == nil {
				return fmt.Errorf("null origin")
			}
			return err
		},
	}
	return x
}

func (this *RpcHandler) MustRegister(c ctx.C, obj any) *RpcHandler {
	err := this.Register(c, obj)
	if err != nil {
		panic(err)
	}
	return this
}

func (this *RpcHandler) Register(c ctx.C, obj any) error {
	var b builder
	err := b.inspect(c, obj)
	if err != nil {
		return err
	}
	this.byPath.Store(b.name, func(c ctx.C, conn *RpcConn, f Frame) error {
		log.Debugf(c, "registered new path %q using %T...", b.name, obj)
		ch := &Channel{
			Conn:   conn,
			byPath: map[string]func(c C, n enc.Node) error{},
			ID:     f.Channel,
		}
		conn.byChan.Store(ch.ID, ch)
		obj := b.build(C{
			C:  c,
			ch: ch,
		}, f.Data)
		log.Debugf(c, "new %+v", obj)
		return nil
	})
	return nil
}
