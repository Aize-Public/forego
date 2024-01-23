package http_test

import (
	"testing"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/http"
	"github.com/Aize-Public/forego/test"
)

func TestCall(t *testing.T) {
	c := test.Context(t)

	s := http.NewServer(c)
	http.CallHandler{
		Path:        "/test/call",
		ReadTimeout: time.Second,
		Handler: func(c ctx.C, call *http.Call) error {
			return http.Error{400, nil}
		},
	}.Register(s)

	addr, err := s.Listen(c, "127.0.0.1:0")
	test.NoError(t, err)

	_, err = http.DefaultClient.Post(c, "http://"+addr.String()+"/test/call", []byte(`[]`))
	test.Error(t, err)
	test.ContainsJSON(c, err.Error(), "Bad Request")
}
