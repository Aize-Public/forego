package ws

import (
	"errors"
	"runtime"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
)

type Channel struct {
	Conn   *RpcConn
	ID     string
	byPath map[string]func(c C, n enc.Node) error
}

// close the channel, removing it from the connection
func (this *Channel) Close(c ctx.C) error {
	log.Infof(c, "closing channel %q", this.ID)
	this.Conn.byChan.Delete(this.ID)
	return this.Conn.Send(c, Frame{
		Channel: this.ID,
		Type:    "close",
	})
}

func (this *Channel) onData(c ctx.C, f Frame) error {
	fn := this.byPath[f.Path]
	if fn == nil {
		return ctx.NewErrorf(c, "no %q for channel %q", f.Path, f.Channel)
	}
	log.Debugf(c, "ch[%q].%q(%v)", f.Channel, f.Path, f.Data)

	ok := func() bool {
		defer func() {
			if r := recover(); r != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				log.Errorf(c, "http/ws: panic serving: %v\n%s", r, buf)
			}
		}()
		err := fn(C{C: c, ch: this}, f.Data)
		if err != nil {
			//log.Warnf(c, "ws: sending %v", err)
			_ = this.Conn.Send(c, Frame{
				Channel: this.ID,
				Type:    "error",
				Data:    enc.MustMarshal(c, err.Error()),
			})
		}
		return true
	}()

	if !ok {
		return errors.New("panic")
	}
	return nil
}

type Frame struct {
	// dialog identifier
	Channel string `json:"channel,omitempty"`

	// routes to an object by path
	Path string `json:"path,omitempty"`

	Type string `json:"type"` // data, error, close

	Data enc.Node `json:"data,omitempty"`
}

type C struct {
	ctx.C
	ch *Channel
}

func (c C) Reply(path string, obj any) error {
	return c.ch.Conn.Send(c, Frame{
		Channel: c.ch.ID,
		Path:    path,
		Data:    enc.MustMarshal(c, obj),
	})
}

// TODO(oha) should we keep it?
//func (c C) Error(obj any) error {
//	return c.ch.Conn.Send(c, Frame{
//		Channel: c.ch.ID,
//		Type:    "error",
//		Data:    enc.MustMarshal(c, obj),
//	})
//}

// Close the websocket
func (c C) Close() error {
	return c.ch.Conn.Close(c, EXIT)
}
