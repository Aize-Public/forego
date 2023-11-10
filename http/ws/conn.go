package ws

import (
	"io"
	"strings"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/utils/sync"
)

type Conn struct {
	h      *Handler
	ws     impl
	byChan sync.Map[string, *Channel]
}

func (this *Conn) Close(c ctx.C, reason int) error {
	return this.ws.Close(c, reason)
}

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
				log.Warnf(c, "inbox: %v", err)
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
			var f Frame
			err := enc.Unmarshal(c, n, &f)
			if err != nil {
				log.Warnf(c, "can't parse request: %v", err)
				continue
			}
			err = this.onData(c, f) // go routine?
			if err != nil {
				log.Errorf(c, "can't process request: %v", err)
				continue
			}
		}
	}
}

func (this *Conn) onData(c ctx.C, f Frame) error {
	switch f.Type {
	case "close":
		// WAIT FOR STUFF?
		return this.Close(c, 1000)
	case "new", "open":
		if h := this.h.byPath.Get(f.Path); h != nil {
			return h(c, this, f)
		}
		return ctx.NewErrorf(c, "unknown path")
	default:
		if ch := this.byChan.Get(f.Channel); ch != nil {
			return ch.onData(c, f)
		}
		return ctx.NewErrorf(c, "unknown channel")
	}
}

func (this *Conn) Send(c ctx.C, f Frame) error {
	return this.ws.Write(c, enc.MustMarshal(c, f))
}

func toLowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[0:1]) + s[1:]
}
