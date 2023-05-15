package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Aize-Public/forego/ctx"
)

// TODO circuit breaker!
type Client struct {
	cli http.Client
}

var DefaultClient = Client{}

func (this Client) Get(c ctx.C, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(c, "GET", url, nil)
	if err != nil {
		return nil, ctx.NewErrorf(c, "can't create request: %w", err)
	}
	res, err := this.Request(req)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
		return io.ReadAll(res.Body)
	}
	return nil, nil
}

func NewRequest(c ctx.C, method string, url string, body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(c, method, url, body)
}

func (this Client) Post(c ctx.C, url string, body []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(c, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, ctx.NewErrorf(c, "can't create request: %w", err)
	}
	res, err := this.Request(req)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
		return io.ReadAll(res.Body)
	}
	return nil, nil
}

func (this Client) Request(r *http.Request) (*http.Response, error) {
	res, err := this.cli.Do(r)
	if err != nil {
		return res, ctx.NewErrorf(r.Context(), "http.Client.Do: %w", err)
	}
	switch res.StatusCode {
	case 200, 204:
		return res, nil
	default:
		return res, Error{res.StatusCode, fmt.Errorf("%s", res.Status)}
	}
}
