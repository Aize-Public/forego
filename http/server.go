package http

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/shutdown"
)

type Server struct {
	mux *http.ServeMux
	h   http.Handler

	// called when a request is done, by default it logs and generate metrics
	OnResponse func(Stat)
}

func NewServer(c ctx.C) *Server {
	this := &Server{
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
			this.mux.ServeHTTP(responseHijacker{w2, w}, r.WithContext(c))
		default:
			this.mux.ServeHTTP(w2, r.WithContext(c))
		}
		this.OnResponse(Stat{
			Method:  r.Method,
			Path:    r.URL.Path,
			UA:      r.UserAgent(),
			Code:    w2.code,
			Elapsed: time.Since(t0),
		})
	})
	this.mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
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
	log.Debugf(c, "listening to %s", addr)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	ch := make(chan error, 1)

	// we start the server in a goroutine
	go func() {
		defer shutdown.Hold().Release()
		// we wrap the listener, so the first call to Accept() will write nil to the error channel
		l := &listener{
			Listener: ln,
			f: func() {
				ch <- nil
			},
		}
		err := s.Serve(l)
		if err != nil {
			log.Warnf(c, "listen(%q) %v", addr, err)
		}
	}()

	go func() {
		<-shutdown.Started()
		c, cf := ctx.WithTimeout(c, 30*time.Second)
		defer cf()
		_ = s.Shutdown(c)
	}()

	// blocks until either an error, or the first Accept() call happen
	return ln.Addr(), <-ch
}

// wrapper for a net listener
type listener struct {
	net.Listener
	once sync.Once
	f    func()
}

func (this *listener) Accept() (net.Conn, error) {
	this.once.Do(this.f)
	return this.Listener.Accept()
}

func (this *Server) HandleRequest(pattern string, f func(c ctx.C, in []byte, r *http.Request) ([]byte, error)) {
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
			log.Warnf(c, "http: %v", err)
			code := ErrorCode(err, 500)
			if code < 500 {
				j, _ := json.Marshal(map[string]any{
					"error":    err.Error(),
					"tracking": nil, // TODO
				})
				_, _ = w.Write(j)
			} else {
				j, _ := json.Marshal(map[string]any{
					// NO CLIENT REPORTING FOR INTERNAL ERRORS (due to security reason, we don't want to leak information about internals)
					"tracking": nil, // TODO
				})
				_, _ = w.Write(j)
			}
			w.WriteHeader(code)
			return
		}

		if len(out) == 0 {
			// no response content
			w.WriteHeader(204)
			return
		}

		w.WriteHeader(200)
		_, err = w.Write(out)
		if err != nil {
			log.Warnf(c, "writing the response: %v", err)
		}
	})
}
