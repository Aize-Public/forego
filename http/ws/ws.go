package ws

import (
	"fmt"
	"io"
	"sync"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
	"golang.org/x/net/websocket"
)

// A request sent from a client to the websocket
type Op[State any] interface {
	Do(ctx.C, State) error
}

// The frame used internally to route the custom data using Path anc Channel
type Frame struct {
	Path     string  `json:"path,omitempty"`
	Channel  string  `json:"channel,omitempty"`
	Data     enc.Map `json:"data,omitempty"`
	Error    string  `json:"error,omitempty"`
	Tracking string  `json:"tracking-id,omitempty"`
}

///////////////////////////////////////////////
// implementation using /x/net/websocket

type ws enc.ReadWriteCloser

type impl struct {
	m    sync.Mutex
	conn *websocket.Conn
}

var _ ws = &impl{}

func (this *impl) Write(c ctx.C, n enc.Node) error {
	this.m.Lock()
	defer this.m.Unlock()
	return this.write(c, enc.JSON{}.Encode(c, n))
}

// note this is not safe for multiple go routines
func (this *impl) write(c ctx.C, data []byte) error {
	w, err := this.conn.NewFrameWriter(websocket.TextFrame)
	if err != nil {
		return ctx.WrapError(c, err)
	}
	ct, err := w.Write(data)
	if err != nil {
		return ctx.WrapError(c, err)
	}
	if ct != len(data) {
		return ctx.NewErrorf(c, "can't write all data: %d < %d", ct, len(data))
	}
	log.Debugf(c, "ws sent: %s", data)
	return nil
}

func (this *impl) Read(c ctx.C) (enc.Node, error) {
	j, err := this.read(c)
	if err != nil {
		return nil, err
	}
	return enc.JSON{}.Decode(c, j)
}

func (this *impl) read(c ctx.C) ([]byte, error) {
	r, err := this.conn.NewFrameReader()
	if err != nil {
		return nil, err
	}
	hr := r.HeaderReader()
	if hr != nil {
		//d, _ := io.ReadAll(hr)
		//log.Debugf(c, "ws head: %x (%q)", d, d)
		_, _ = io.Copy(io.Discard, hr) // we don't care
	}
	if c.Err() != nil {
		return nil, ctx.WrapError(c, c.Err())
	}

	switch r.PayloadType() {
	case websocket.ContinuationFrame:
		return nil, ctx.NewErrorf(c, "unsupported continuation")
	case websocket.TextFrame, websocket.BinaryFrame:
	case websocket.CloseFrame:
		// TODO read message and debug?
		return nil, ctx.WrapError(c, io.EOF)
	case websocket.PingFrame, websocket.PongFrame:
		panic("unsupported ping-pong")
		/*
			            // TODO implements this
						b := make([]byte, maxControlFramePayloadLength)
						n, err := io.ReadFull(frame, b)
						if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
							return nil, err
						}
						io.Copy(ioutil.Discard, frame)
						if frame.PayloadType() == PingFrame {
							if _, err := handler.WritePong(b[:n]); err != nil {
								return nil, err
							}
						}
						return nil, nil
		*/
	default:
		panic(fmt.Sprintf("unexpected type %b", r.PayloadType()))
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, ctx.WrapError(c, err)
	}
	log.Debugf(c, "ws recv: %s", data)
	return data, nil
}

func (this *impl) Close() error {
	return this.conn.Close()
}
