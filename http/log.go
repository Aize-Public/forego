package http

import (
	"bufio"
	"net"
	"net/http"
	"time"
)

type Stat struct {
	Method  string
	Path    string
	UA      string
	Code    int
	Elapsed time.Duration
}

/*
func defaultMiddleware(h http.Handler, f func(s Stat)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		c := r.Context()
		c = ctx.WithTag(c, "ua", r.UserAgent())
		c = ctx.WithTag(c, "path", r.URL.Path)
		w2 := &response{w, 0}
		switch w := w.(type) {
		case http.Hijacker:
			h.ServeHTTP(responseHijacker{w2, w}, r)
		default:
			h.ServeHTTP(w2, r)
		}
		f(Stat{
			Path:    r.URL.Path,
			UA:      r.UserAgent(),
			Code:    w2.code,
			Elapsed: time.Since(t0),
		})
	})
}
*/

type response struct {
	http.ResponseWriter
	code int
}

type responseHijacker struct {
	*response
	hijacker http.Hijacker
}

var _ http.Hijacker = &responseHijacker{}

func (r *response) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *response) Write(b []byte) (int, error) {
	if r.code == 0 {
		r.code = 200
	}
	return r.ResponseWriter.Write(b)
}

func (r responseHijacker) Hijack() (conn net.Conn, rw *bufio.ReadWriter, err error) {
	return r.hijacker.Hijack()
}
