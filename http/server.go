package http

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

type Server struct {
	mux *http.ServeMux
	h   http.Handler

	// called when a request is done, by default it logs and generate metrics
	OnResponse func(Stat)
}

func NewServer(c ctx.C, cb func(r Stat)) Server {
	this := Server{
		mux: http.NewServeMux(),
		OnResponse: func(r Stat) {
			log.Infof(c, "%s %d in %v", r.Path, r.Code, r.Elapsed)
			// TODO metrics
		},
	}
	this.h = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		c := r.Context()
		c = ctx.WithTag(c, "ua", r.UserAgent())
		c = ctx.WithTag(c, "path", r.URL.Path)

		w2 := &response{w, 0}
		switch w := w.(type) {
		case http.Hijacker:
			// if there is an hijacker, we need to be a bit clever
			this.mux.ServeHTTP(responseHijacker{w2, w}, r)
		default:
			this.mux.ServeHTTP(w2, r)
		}

		this.OnResponse(Stat{
			Method:  r.Method,
			Path:    r.URL.Path,
			UA:      r.UserAgent(),
			Code:    w2.code,
			Elapsed: time.Since(t0),
		})
	})
	return this
}

func (this Server) Mux() *http.ServeMux {
	return this.mux
}

func (this Server) Listen(c ctx.C, addr string) (net.Addr, error) {
	s := http.Server{
		BaseContext: func(l net.Listener) context.Context {
			return ctx.WithTag(c, "http.addr", l.Addr().String())
		},
		ConnContext: func(c context.Context, conn net.Conn) context.Context {
			return ctx.WithTag(c, "http.remote", conn.RemoteAddr().String())
		},
		ConnState: func(conn net.Conn, state http.ConnState) {
			// TODO(oha): do we need to setup a limiter? if so, this is to know when any hijacker kicks in
		},
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      this.h,
	}
	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return ln.Addr(), err
	}

	// we create a channel for the error
	ch := make(chan error)

	// we start the server in a goroutine
	go func() {
		err := s.Serve(ln)
		select {
		case ch <- err:
		default:
			if c.Err() == nil {
				log.Errorf(c, "http closing: %v", err)
			}
		}
	}()

	// if it doesn't fails in a second, we consider it good to go!
	select {
	case err := <-ch:
		return ln.Addr(), err
	case <-time.After(time.Second):
		return ln.Addr(), nil
	}
}

func (this *Server) OnRequest(pattern string, f func(c ctx.C, in []byte, r *http.Request) ([]byte, error)) {
	this.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		out, err := func() ([]byte, error) {
			var in []byte
			var err error
			if r.Body != nil {
				in, err = io.ReadAll(r.Body)
				if err != nil {
					return nil, NewErrorf(c, 400, "can't read body: %w", err)
				}
			}
			return f(c, in, r)
		}()
		if err != nil {
			w.WriteHeader(ErrorCode(err, 500))
			return
		}
		if len(out) == 0 {
			w.WriteHeader(204)
			return
		}
		w.WriteHeader(200)
		_, err = w.Write(out)
		if err != nil {
			log.Warnf(c, "writing the reponse: %v", err)
		}
	})
}
