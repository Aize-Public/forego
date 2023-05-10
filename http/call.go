package http

import (
	"net/http"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/utils"
)

type CallHandler struct {
	Path           string
	MaxLength      int           // default is 1Mb
	RequestTimeout time.Duration // default is 30s
	Handler        func(c ctx.C, call *Call) error
}

func (this CallHandler) requestTimeout() time.Duration {
	if this.RequestTimeout > 0 {
		return this.RequestTimeout
	}
	return 30 * time.Second
}

func (this CallHandler) Register(s *Server) {
	s.mux.HandleFunc(this.Path, func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		out, err := func() ([]byte, error) {
			c, cf := ctx.WithTimeout(c, this.requestTimeout())
			defer cf()
			call := Call{
				r: r,
				w: w,
			}
			var err error
			if r.Body != nil {
				call.reqBody, err = utils.ReadAll(c, r.Body)
				if err != nil {
					return nil, NewErrorf(c, 400, "can't read body: %w", err)
				}
			}
			err = this.Handler(c, &call)
			return call.resBody, err
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

type Call struct {
	r       *http.Request
	reqBody []byte

	w       http.ResponseWriter
	resBody []byte
}
