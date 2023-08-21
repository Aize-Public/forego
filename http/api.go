package http

import (
	"bytes"
	"net/http"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
)

type Doable interface {
	Do(ctx.C) error
}

func (s *Server) MustRegisterAPI(c ctx.C, obj Doable) {
	err := s.RegisterAPI(c, obj)
	if err != nil {
		panic(err)
	}
}

func (s *Server) RegisterAPI(c ctx.C, obj Doable) error {
	handler, err := api.NewServer(c, obj)
	if err != nil {
		return err
	}
	f := func(c ctx.C, in []byte, r *http.Request) ([]byte, error) {
		req := &api.JSON{}
		if r.Body != nil {
			err := req.ReadFrom(c, bytes.NewBuffer(in))
			if err != nil {
				return nil, ctx.NewErrorf(c, "can't read request body: %v", err)
			}
			defer r.Body.Close()
		} else {
			log.Infof(c, "can/t get body: %v", err)
		}
		// TODO auth

		obj, err := handler.Recv(c, req)
		if err != nil {
			return nil, err
		}
		err = obj.Do(c)
		if err != nil {
			return nil, err
		}
		log.Debugf(c, "API %#v", obj)

		res := &api.JSON{}
		err = handler.Send(c, obj, res)
		if err != nil {
			return nil, err
		}
		return enc.JSON{}.Encode(c, res.Data), nil
	}

	urls := handler.URLs()
	if len(urls) == 0 {
		return ctx.NewErrorf(c, "no URL to register for %T", obj)
	}

	for _, u := range handler.URLs() {
		log.Debugf(c, "registering to %q", u.Path)
		s.HandleRequest(u.Path, f)
	}
	return nil
}
