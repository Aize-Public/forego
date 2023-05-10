package api_test

import (
	"strings"
	"testing"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/test"
)

func TestAPI(t *testing.T) {
	c := test.C(t)

	h, err := api.NewHandler(c, WordCount{})
	test.NoError(t, err)

	req := api.JSON{}
	obj, err := h.Server().Recv(c, req)
	test.NoError(t, err)
	t.Logf("obj: %+v", obj)

	err = obj.Foo(c)
	test.NoError(t, err)

	res := api.JSON{}
	err = h.Server().Send(c, obj, res)
	test.NoError(t, err)
	t.Logf("res: %s", res)
}

type UID string

type WordCount struct {
	UID UID    `api:"auth,required"`
	Str string `api:"in" json:"str"`
	Ct  int    `api:"out" json:"ct"`
}

func (this *WordCount) Foo(ctx.C) error {
	this.Ct = len(strings.Split(this.Str, " "))
	return nil
}
