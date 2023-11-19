package http_test

import (
	"bytes"
	gohttp "net/http"
	"net/url"
	"testing"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/http"
	"github.com/Aize-Public/forego/test"
)

type UID string

type Inc struct {
	R   api.Request `url:"/inc"`
	UID UID         `api:"auth"`

	Name    string `api:"in,out" json:"name"`
	Amount  int    `api:"in" json:"amount"`
	Current int    `api:"out" json:"current"`

	State map[string]int
}

func (this *Inc) Do(c ctx.C) error {
	log.Debugf(c, "%+v.Do()", *this)
	if this.Amount < 0 {
		return http.NewErrorf(c, 400, "amount can't be negative")
	}
	this.Current = this.State[this.Name] + this.Amount
	this.State[this.Name] = this.Current
	return nil
}

func TestAPI(t *testing.T) {
	c := test.Context(t)

	s := http.NewServer(c)
	_, err := s.RegisterAPI(c, &Inc{
		State: map[string]int{},
	})
	test.NoError(t, err)

	{
		req, err := http.NewRequest(c, "POST", "/inc", bytes.NewBufferString(`{"name":"foo","amount":3}`))
		test.NoError(t, err)
		w := &ResponseWriter{}
		s.Mux().ServeHTTP(w, req)
		res := w.Buf.String()
		t.Logf("res: %d %s", w.Code, res)
		test.EqualsGo(t, 200, w.Code)
		test.Contains(t, res, `:3,`)
	}

	// we must use the network to test the client
	addr, err := s.Listen(c, "127.0.0.1:0")
	test.NoError(t, err)

	cli := http.Client{
		BaseUrl: &url.URL{Scheme: "http", Host: addr.String()}, // this will be used to build relative urls from
	}

	{
		op := Inc{
			Name:   "foo",
			Amount: 42,
		}
		err := cli.API(c, &op)
		test.NoError(t, err)
		test.EqualsGo(t, 45, op.Current)
	}
}

// helpers

type ResponseWriter struct {
	Buf    bytes.Buffer
	header gohttp.Header
	Code   int
}

var _ gohttp.ResponseWriter = &ResponseWriter{}

func (this *ResponseWriter) Header() gohttp.Header {
	if this.header == nil {
		this.header = gohttp.Header{}
	}
	return this.header
}

func (this *ResponseWriter) Write(b []byte) (int, error) {
	if this.Code == 0 {
		this.WriteHeader(200)
	}
	return this.Buf.Write(b)
}

func (this *ResponseWriter) WriteHeader(code int) {
	this.Code = code
}
