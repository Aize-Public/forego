package ws

import (
	"reflect"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/enc"
	"github.com/google/uuid"
)

type Client[State any] struct {
	// generic function for sending data to the server
	send func(ctx.C, enc.Node) error

	// cache for api.Client
	apiType map[reflect.Type]api.Client[Op[State]]
	apiPath map[string]api.Client[Op[State]]

	// handlers by channel
	handlers map[string]func(ctx.C, Frame) error

	// fallback for anything that has no handler
	fallback func(ctx.C, enc.Node) error
}

// create a client for this websocket, you must provide a way to send data to the server, and the fallback for any messages which
// are not a reply to our requests
func (h Handler[State]) NewClient(c ctx.C, send func(ctx.C, enc.Node) error, fallback func(ctx.C, enc.Node) error) *Client[State] {
	this := &Client[State]{
		send:     send,
		apiType:  map[reflect.Type]api.Client[Op[State]]{},
		apiPath:  map[string]api.Client[Op[State]]{},
		handlers: map[string]func(ctx.C, Frame) error{},
		fallback: func(c ctx.C, n enc.Node) error {
			return ctx.NewErrorf(c, "unexpected %+v", n)
		},
	}
	for _, h := range h.resolver {
		hc := h.Client()
		this.apiType[h.Type()] = hc
		for _, path := range h.Paths() {
			this.apiPath[path] = hc
		}
	}
	return this
}

func (this *Client[State]) onRecv(c ctx.C, n enc.Node) error {
	var f Frame
	err := enc.Unmarshal(c, n, &f)
	if err != nil {
		return this.fallback(c, n)
	}
	h, exists := this.handlers[f.Channel]
	if !exists {
		return this.fallback(c, n)
	}
	return h(c, f)
}

// send the request to the server on a new channel, and use the given callback function for any reply
// returns an object which can then be closed
func (this *Client[State]) Request(c ctx.C, req Op[State], fn func(ctx.C, Op[State]) error) (*cliRequest[State], error) {
	h, err := api.NewClient(c, req)
	if err != nil {
		return nil, err
	}
	r := &cliRequest[State]{
		h:   h,
		ch:  uuid.NewString(),
		cli: this,
	}
	this.handlers[r.ch] = func(c ctx.C, f Frame) error {
		h, exists := this.apiPath[f.Path]
		if !exists {
			return ctx.NewErrorf(c, "unexpected path %q", f.Path)
		}
		j := &api.JSON{
			Data: f.Data.(enc.Map), // TODO FIXME
		}
		err := h.Client().Recv(c, j, req)
		if err != nil {
			return err
		}
		return fn(c, req)
	}
	return r, r.send(c, req)
}

type cliRequest[State any] struct {
	h   api.Client[Op[State]]
	ch  string
	cli *Client[State]
}

func (this cliRequest[State]) Close() error {
	delete(this.cli.handlers, this.ch)
	return nil
}

func (this cliRequest[State]) send(c ctx.C, op Op[State]) error {
	req := api.JSON{}
	err := this.h.Send(c, op, &req)
	if err != nil {
		return err
	}
	f := Frame{
		Path:    this.h.Path(),
		Data:    req.Data,
		Channel: this.ch,
	}
	return this.cli.send(c, enc.MustMarshal(c, f))
}
