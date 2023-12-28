package http

import (
	"bufio"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Aize-Public/forego/api/openapi"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/utils"
)

func (this *Server) HandleFunc(path string, h http.HandlerFunc) *openapi.Path {
	return this.Handle(path, http.HandlerFunc(h))
}

func (this *Server) Handle(path string, h http.Handler) *openapi.Path {
	this.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()
		c := r.Context()
		c = ctx.WithTag(c, "ua", r.UserAgent())
		c = ctx.WithTag(c, "path", r.URL.Path)

		w2 := &response{w, 0}
		switch w := w.(type) {
		case http.Hijacker:
			// if there is an hijacker, we need to be a bit clever
			h.ServeHTTP(responseHijacker{w2, w}, r.WithContext(c))
		default:
			h.ServeHTTP(w2, r.WithContext(c))
		}

		metric{
			Method: r.Method,
			Code:   w2.code,
			Path:   path,
		}.observe(time.Since(t0))
	})

	p := &openapi.Path{}
	this.OpenAPI.Paths[path] = p
	return p
}

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
	if r.code != 0 {
		stack := utils.Stack(1, 10)
		log.Warnf(nil, "duplicate WriteHeader() at %s", strings.Join(stack, "\n"))
		return
	}
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
