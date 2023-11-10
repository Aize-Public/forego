package ws

import (
	"net/http"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/shutdown"
	"github.com/Aize-Public/forego/utils/sync"
	"golang.org/x/net/websocket"
)

type Handler struct {
	byPath sync.Map[string, func(ctx.C, *Conn, Frame) error]
}

// return a websocket.Server which can be used as an http.Handler
func (this *Handler) Server() websocket.Server {
	x := websocket.Server{
		Handler: websocket.Handler(func(conn *websocket.Conn) {
			c := conn.Request().Context()
			c, cf := ctx.Span(c, "ws")
			defer cf(nil)

			defer shutdown.Hold().Release()

			//defer metrics.WS{Path: path}.Start().End(c)
			ws := Conn{
				h: this,
				ws: &wsImpl{
					conn: conn,
				},
			}
			defer ws.Close(c, 1000)
			ws.Loop(c)
		}),
		Handshake: func(c *websocket.Config, r *http.Request) error {
			// NOTE(oha): since the current WS implementation require a jwt token to be sent as a payload
			// there is no risk of cross domain forgery, since they will still need to grab the token in the first place
			//log.Debugf(r.Context(), "WS: handshake origin: %q", r.Header.Get("Origin"))
			return nil
		},
	}
	return x
}

func (this *Handler) Register(c ctx.C, obj any) error {
	var b builder
	err := b.inspect(c, obj)
	if err != nil {
		return err
	}
	this.byPath.Store(b.name, func(c ctx.C, conn *Conn, f Frame) error {
		log.Debugf(c, "new %q...", b.name)
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
