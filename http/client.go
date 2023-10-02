package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

// TODO circuit breaker!
type Client struct {
	cli     http.Client
	BaseUrl *url.URL
	// TODO do we need a "proxy" host which will be used instead of the url host?
}

var DefaultClient = Client{}

func (this Client) Get(c ctx.C, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(c, "GET", url, nil)
	if err != nil {
		return nil, ctx.NewErrorf(c, "can't create request: %w", err)
	}
	res, err := this.Do(req)
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
	req.Header.Set(`Content-Type`, `application/json`)
	res, err := this.Do(req)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
		return io.ReadAll(res.Body)
	}
	return nil, nil
}

func (this Client) Do(r *http.Request) (*http.Response, error) {
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

func (this Client) API(c ctx.C, obj Doable) error {
	c = ctx.WithTag(c, "api", fmt.Sprintf("client %T", obj))
	h, err := api.NewClient(c, obj)
	if err != nil {
		return err
	}
	{
		data := &api.JSON{}
		err = h.Send(c, obj, data)
		if err != nil {
			return err
		}
		j, err := json.Marshal(data.Data)
		if err != nil {
			return ctx.NewErrorf(c, "can't marshal %T: %w", obj, err)
		}

		u := h.URL()
		if this.BaseUrl != nil {
			u, err = this.BaseUrl.Parse(u.String())
			if err != nil {
				return ctx.NewErrorf(c, "can't rebase url: %w", err)
			}
		}

		req, err := NewRequest(c, "POST", u.String(), bytes.NewBuffer(j))
		if err != nil {
			return err
		}

		log.Debugf(c, "client[%T].Send() %s", obj, j)
		res, err := this.Do(req)
		if err != nil {
			return ctx.NewErrorf(c, "can't send request: %w", err)
		}
		switch res.StatusCode {
		case 204:
			log.Debugf(c, "204 no response")
			return nil
		case 200:
			err := data.ReadFrom(c, res.Body)
			if err != nil {
				return ctx.NewErrorf(c, "can't read response: %w", err)
			}
			res.Body.Close()
			log.Debugf(c, "client[%T].Recv() %s", obj, j)
			return h.Recv(c, data, obj)
		default:
			return ctx.NewErrorf(c, "can't connect: %s", res.Status)
		}
	}
}
