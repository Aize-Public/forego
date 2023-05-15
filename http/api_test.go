package http_test

import (
	"bytes"
	gohttp "net/http"
	"net/url"
	"testing"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/http"
	"github.com/Aize-Public/forego/test"
)

type Inc struct {
	R       api.Request `url:"/inc"`
	Name    string      `api:"in,out" json:"name"`
	Amount  int         `api:"in" json:"amount"`
	Current int         `api:"out" json:"current"`

	State map[string]int
}

func (this *Inc) Do(c ctx.C) error {
	if this.Amount < 0 {
		return http.NewErrorf(c, 400, "amount can't be negative")
	}
	this.Current = this.State[this.Name] + this.Amount
	this.State[this.Name] = this.Current
	return nil
}

func TestAPI(t *testing.T) {
	c := test.C(t)

	s := http.NewServer(c)
	err := http.API(c, s, &Inc{
		State: map[string]int{},
	})
	test.NoError(t, err)

	{
		req, err := http.NewRequest(c, "POST", "/inc", bytes.NewBufferString(`{"name":"foo","amount":3}`))
		test.NoError(t, err)
		w := &ResponseWriter{}
		s.Mux().ServeHTTP(w, req)
		res := w.Buf.String()
		test.Assert(t, w.Code == 200)
		t.Logf("res: %s", res)
		test.ContainsJSON(t, res, `:3,`)
	}

	addr, err := s.Listen(c, "127.0.0.1:0")
	test.NoError(t, err)

	cli := http.Client{
		BaseUrl: &url.URL{Scheme: "http", Host: addr.String()},
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
	return this.Buf.Write(b)
}

func (this *ResponseWriter) WriteHeader(code int) {
	this.Code = code
}
