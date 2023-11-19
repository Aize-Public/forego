package http

import (
	"net/http"

	"github.com/Aize-Public/forego/api/openapi"
)

func (this *Server) Handle(path string, h http.Handler) *openapi.Path {
	this.mux.Handle(path, h)
	p := &openapi.Path{}
	this.OpenAPI.Paths[path] = p
	return p
}

func (this *Server) HandleFunc(path string, h http.HandlerFunc) *openapi.Path {
	this.mux.HandleFunc(path, h)
	p := &openapi.Path{}
	this.OpenAPI.Paths[path] = p
	return p
}
