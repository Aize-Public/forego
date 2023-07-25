package ws

import (
	gohttp "net/http"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/shutdown"
	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

/*
Handler is used to define a websocket.
It tells on which path to mount it, what to do on connect and which object should be expected to be processed
*/
type Handler[State any] struct {
	// optional callback fired on connection. No error can be returned, but you are free to close the connection if needed
	OnConnect func(c ctx.C, conn *Conn[State])

	// fired when a shutdown start, if empty it will send `"shutdown"` to the client
	OnShutdown func(c ctx.C, conn *Conn[State])

	// fire after close
	OnExit func(c ctx.C, conn *Conn[State])

	resolver map[string]api.Handler[Op[State]] // TODO(oha) change to something specific to websocket, e.g. that allows parallelism or other flows
	impl     websocket.Server
}

var _ gohttp.Handler = (*Handler[any])(nil)

// Build a new Websocket handler, which support the given operations.
// Note: the given Op objects might not have an `api:"url"` field, in which case the Object name will be used
func New[State any](c ctx.C, objs ...Op[State]) (*Handler[State], error) {
	this := &Handler[State]{
		OnShutdown: func(c ctx.C, conn *Conn[State]) {
			_ = conn.Send(c, "shutdown")
		},
	}

	this.resolver = map[string]api.Handler[Op[State]]{}
	for _, obj := range objs {
		s, err := api.NewHandler(c, obj)
		if err != nil {
			return nil, err
		}
		for _, p := range s.Paths() {
			this.resolver[p] = s
		}
	}
	this.impl = websocket.Server{
		Handler: websocket.Handler(func(conn *websocket.Conn) {
			c := conn.Request().Context()
			ws := this.connect(c, &impl{conn: conn})
			defer ws.Close(c, "exit")
			ws.loop(c)
		}),
	}
	return this, nil
}

func (this *Handler[State]) connect(c ctx.C, ws ws) *Conn[State] {
	c, cf := ctx.Span(c, "ws")
	defer cf(nil)

	sid := uuid.NewString()
	c = ctx.WithTag(c, "sid", sid)

	defer shutdown.Hold().Release() // prevents shutdown

	//defer metrics.WS{Path: path}.Start().End(c)

	conn := &Conn[State]{
		h:   this,
		sid: sid,
		ws:  ws,
	}

	if this.OnConnect != nil {
		this.OnConnect(c, conn)
	}
	return conn
}

func (this *Handler[State]) ServeHTTP(w gohttp.ResponseWriter, r *gohttp.Request) {
	this.impl.ServeHTTP(w, r)
}
