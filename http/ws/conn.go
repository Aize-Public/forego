package ws

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/shutdown"
	"github.com/Aize-Public/forego/utils/lists"
)

type Conn[State any] struct {
	h      *Handler[State]
	m      sync.Mutex
	sid    string
	UID    enc.Node
	closed bool
	ws     ws

	State State

	onClose []func()
}

func (this *Conn[State]) SID() string {
	return this.sid
}

func (this *Conn[State]) OnClose(c ctx.C, f func()) {
	this.m.Lock()
	defer this.m.Unlock()
	this.onClose = append(this.onClose, f)
}

// Send a generic payload (no frame is added, which means the client won't recognize this as a reply to a channel or anything like that)
// most likely you don't want to use this
func (this *Conn[State]) Send(c ctx.C, obj any) error {
	this.m.Lock()
	defer this.m.Unlock()
	if this.closed {
		return ctx.NewErrorf(c, "closed")
	}
	n, err := enc.Marshal(c, obj)
	if err != nil {
		return err
	}
	return this.ws.Write(c, n)
}

func (this *Conn[State]) loop(c ctx.C) {
	// To make thing safe and not leaking, we follow this pattern
	// no matter what, if the connectsion closes, the read-loop closes as well
	// which then closes the inbox
	// which then closes the operational loop
	// which then returns
	// this means we should always try to close the connect to abort
	inbox := make(chan enc.Node)
	go func() {
		defer close(inbox) // this will tear all down
		for {
			n, err := this.ws.Read(c) //NOTE(oha): read(c) might not be able to honor c.Done()
			if err != nil {
				if errors.Is(err, io.EOF) {
					//log.Debugf(c, "ws read: %v", err)
				} else {
					log.Warnf(c, "can't ws read: %v (%v)", err, c.Err())
				}
				return
			}
			select {
			case inbox <- n:
			case <-c.Done():
				if err != nil {
					log.Debugf(c, "ws read loop aborted: %v", c.Err())
					return
				}
			}
		}
	}()

	shutwarn := shutdown.Started()
	for {
		select {
		case <-c.Done():
			log.Warnf(c, "cancel: %v", ctx.Cause(c))
			this.Close(c, "cancel")
			return

		case <-shutwarn:
			this.h.OnShutdown(c, this)
			shutwarn = nil // don't fire again
			time.AfterFunc(5*time.Second, func() {
				this.Close(c, "shutdown")
			})

		case n, ok := <-inbox:
			if !ok {
				log.Debugf(c, "inbox closed")
				return
			}
			c = ctx.WithTracking(c, "")
			err := this.onData(c, n)
			if err != nil {
				log.Warnf(c, "can't understand the client: %v", err)
				err = this.Send(c, Frame{
					Error:    err.Error(), // TODO only for 4xx? close on others?
					Tracking: ctx.GetTracking(c),
				})
				if err != nil {
					log.Warnf(c, "can't write: %v", err)
					return
				}
			}
		}
	}
}

func (this *Conn[State]) onData(c ctx.C, n enc.Node) error {
	var f Frame
	err := enc.Unmarshal(c, n, &f)
	if err != nil {
		return err
	}
	log.Debugf(c, "onData(%+v) => %+v", n, f)
	s, found := this.h.resolver[f.Path]
	if !found {
		return ctx.NewErrorf(c, "invalid path %q", f.Path)
	}
	req := api.JSON{
		Data: f.Data.(enc.Map), // TODO FIXME
		UID:  this.UID,
	}
	obj, err := s.Server().Recv(c, &req)
	if err != nil {
		return err
	}
	r := Request[State]{
		Conn:    this,
		Channel: f.Channel,
	}
	err = obj.Do(c, r)
	if err != nil {
		_ = r.Error(c, err)
		return nil
	}
	res := api.JSON{}
	err = s.Server().Send(c, obj, &res)
	if err != nil {
		return err
	}
	return this.Send(c, Frame{
		Path:     s.Path(),
		Channel:  f.Channel,
		Data:     res.Data,
		Tracking: ctx.GetTracking(c),
	})
}

func (this *Conn[State]) Close(c ctx.C, reason string) error {
	this.m.Lock()
	defer this.m.Unlock()
	if this.closed {
		return ctx.NewErrorf(c, "dup close: %q", reason)
	}
	fs := this.onClose
	lists.Reverse(fs)
	for _, f := range fs {
		f()
	}
	// TODO send reason
	err := this.ws.Close() // code 1000
	return ctx.WrapError(c, err)
}

type Request[State any] struct {
	*Conn[State]
	Channel string
}

// send a reply to the same channel
func (this Request[State]) Reply(c ctx.C, obj any) error {
	n, err := enc.Marshal(c, obj)
	if err != nil {
		return err
	}
	return this.Conn.Send(c, Frame{
		Channel: this.Channel,
		Data:    n, // TODO(oha): can we make this generic?
	})
}

// send an application error to the client on the same channel
func (this Request[State]) Error(c ctx.C, err error) error {
	return this.Conn.Send(c, Frame{
		Channel:  this.Channel,
		Error:    err.Error(),
		Tracking: ctx.GetTracking(c),
	})
}
