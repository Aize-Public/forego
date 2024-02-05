package ws

import (
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/utils/sync"
)

// extends a websocket connection
// binds different objects instances to different channels, and allow for RPC on them
// it uses a `Frame` to open channels and route messages to them
type RpcConn struct {
	Conn
	h      *RpcHandler
	byChan sync.Map[string, *Channel]
}

func (this *RpcHandler) newRpcConn(impl impl) *RpcConn {
	conn := &RpcConn{
		Conn: Conn{
			ws: impl,
		},
		h: this,
	}
	conn.Conn.onData = func(c ctx.C, n enc.Node) error {
		var f Frame
		err := enc.Unmarshal(c, n, &f)
		if err != nil {
			return err
		}
		return conn.onFrame(c, f)
	}
	return conn
}

func (this *RpcConn) onFrame(c ctx.C, f Frame) error {
	//var f Frame
	//err := enc.Unmarshal(c, n, &f)
	//if err != nil {
	//	return ctx.NewErrorf(c, "can't parse request: %v", err)
	//}

	switch f.Type {
	case "close":
		// WAIT FOR STUFF?
		return this.Close(c, 1000)
	case "new", "open":
		if fn := this.h.byPath.Get(f.Path); fn != nil {
			return fn(c, this, f)
		}
		return ctx.NewErrorf(c, "unknown path: %q", f.Path)
	default:
		if ch := this.byChan.Get(f.Channel); ch != nil {
			return ch.onData(c, f)
		}
		return ctx.NewErrorf(c, "unknown channel")
	}
}
