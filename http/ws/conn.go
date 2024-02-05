package ws

import (
	"io"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
)

// a server side generic websocket connection
type Conn struct {
	ws  impl
	sid string
	//byChan sync.Map[string, *Channel]
	onData func(c ctx.C, n enc.Node) error
}

func (this Conn) SID() string { return this.sid }

func (this *Conn) Close(c ctx.C, reason int) error {
	return this.ws.Close(c, reason)
}

// main loop, which sends evento to the connection HandlerFunc
func (this *Conn) Loop(c ctx.C) error {
	inbox := make(chan enc.Node)
	go func() {
		defer close(inbox)
		for {
			n, err := this.ws.Read(c)
			switch err {
			case io.EOF:
				log.Debugf(c, "inbox: EOF")
				return
			default:
				if c.Err() == nil { // ignore cancels
					log.Warnf(c, "inbox: %v", err)
				}
				return
			case nil:
				select {
				case inbox <- n:
				case <-c.Done():
					log.Warnf(c, "inbox: %v", err)
					return
				}
			}
		}
	}()
	defer this.Close(c, 1000)
	for {
		select {
		case <-c.Done():
			return c.Err()
		case n, ok := <-inbox:
			if !ok {
				return ctx.NewErrorf(c, "inbox closed")
			}
			err := this.onData(c, n) // we want no implicit parallelism here
			if err != nil {
				this.Send(c, Frame{
					Type: "disconnect",
					Data: enc.MustMarshal(c, err),
				})
				return ctx.NewErrorf(c, "closing: %v", err)
			}
		}
	}
}

func (this *Conn) Send(c ctx.C, obj any) error {
	return this.ws.Write(c, enc.MustMarshal(c, obj))
}
