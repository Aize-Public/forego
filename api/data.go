package api

import (
	"reflect"

	"github.com/Aize-Public/forego/ctx"
)

// the client request object used to Marshal a request to a server
type ClientRequest interface {
	Marshal(c ctx.C, name string, into reflect.Value) error
}

// the client response object used to unmarshal the response from the server
type ClientResponse interface {
	Unmarshal(c ctx.C, name string, from reflect.Value) error
}

// the server request object used to Unmarshal the request from the client
type ServerRequest interface {
	Unmarshal(c ctx.C, name string, into reflect.Value) error
	Auth(c ctx.C, into reflect.Value, required bool) error
}

// the server response object used to marshal the response to the client
type ServerResponse interface {
	Marshal(c ctx.C, name string, into reflect.Value) error
}
