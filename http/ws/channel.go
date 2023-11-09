package ws

import (
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
)

type Channel struct {
	Conn   *Conn
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
	h := this.byPath[f.Path]
	if h == nil {
		return ctx.NewErrorf(c, "no %q for channel %q", f.Path, f.Channel)
	}
	log.Debugf(c, "ch[%q].%q(%v)", f.Channel, f.Path, f.Data)
	err := h(C{
		C:  c,
		ch: this,
	}, f.Data)
	if err != nil {
		this.Conn.Send(c, Frame{
			Channel: this.ID,
			Type:    "error",
			Data:    enc.MustMarshal(c, err.Error()),
		})
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

func (c C) Error(obj any) error {
	return c.ch.Conn.Send(c, Frame{
		Channel: c.ch.ID,
		Type:    "error",
		Data:    enc.MustMarshal(c, obj),
	})
}

func (c C) Close() error {
	return c.ch.Close(c)
}
