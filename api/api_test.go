package api_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/test"
)

func TestAPI(t *testing.T) {
	c := test.C(t)
	alice := UID("alice")

	t.Logf("handler...")
	h, err := api.NewHandler(c, WordCount{})
	test.NoError(t, err)

	obj := WordCount{
		Str: "foo bar",
	}
	data := &api.JSON{}

	t.Logf("client send...")
	err = h.Client().Send(c, obj, data)
	test.NoError(t, err)
	t.Logf("request: %s", data.Data)
	test.NotEmpty(t, data.Data)

	t.Logf("auth...")
	data.UID, _ = json.Marshal(alice)

	{
		t.Logf("server recv...")
		obj, err := h.Server().Recv(c, data)
		test.NoError(t, err)
		test.NotEmpty(t, obj)
		test.EqualsJSON(t, alice, obj.UID)

		t.Logf("Foo()ing...")
		err = obj.Foo(c)
		test.NoError(t, err)

		t.Logf("server send...")
		data = &api.JSON{} // clean up
		err = h.Server().Send(c, obj, data)
		test.NoError(t, err)
		t.Logf("response: %s", data.Data)
	}

	t.Logf("client recv...")
	err = h.Client().Recv(c, data, &obj)
	test.NoError(t, err)

	test.Assert(t, obj.Ct == 2)
}

type UID string

type WordCount struct {
	UID UID    `api:"auth,required" json:"uid"`
	Str string `api:"in" json:"str"`
	Ct  int    `api:"out" json:"ct"`
}

func (this *WordCount) Foo(ctx.C) error {
	this.Ct = len(strings.Split(this.Str, " "))
	return nil
}
