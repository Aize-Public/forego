package api_test

import (
	"strings"
	"testing"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

func TestRequired(t *testing.T) {
	c := test.Context(t)

	t.Logf("handler...")
	h, err := api.NewHandler(c, &WordCount{})
	test.NoError(t, err)

	ser := h.Server()

	data := &api.JSON{}
	data.Data = enc.Map{} // empty request should give a 400
	data.UID, _ = enc.Marshal(c, UID("alice"))

	t.Logf("server recv...")
	_, err = ser.Recv(c, data)
	test.Error(t, err)
}

func TestAPI(t *testing.T) {
	c := test.Context(t)
	alice := UID("alice")

	t.Logf("handler...")
	h, err := api.NewHandler(c, &WordCount{})
	test.NoError(t, err)

	ser := h.Server()
	cli := h.Client()

	obj := WordCount{
		Str: "foo bar",
	}
	data := &api.JSON{}

	t.Logf("client send...")
	err = cli.Send(c, &obj, data)
	test.NoError(t, err)

	t.Logf("request: %s", data.Data)
	test.NotEmpty(t, data.Data)

	t.Logf("auth...")
	// as example, we just inject a string here, implementation should use
	// jwt tokens or other authentication mechanism, and place the
	// encoded result in the .UID field here
	data.UID, _ = enc.Marshal(c, alice)

	{
		t.Logf("server recv...")
		obj, err := ser.Recv(c, data)
		test.NoError(t, err)
		test.NotEmpty(t, obj)
		test.EqualsJSON(c, alice, obj.UID)

		t.Logf("Foo()ing...")
		err = obj.Do(c)
		test.NoError(t, err)

		t.Logf("server send...")
		data = &api.JSON{} // clean up
		err = ser.Send(c, obj, data)
		test.NoError(t, err)
		t.Logf("response: %s", data.Data)
	}

	t.Logf("client recv...")
	err = cli.Recv(c, data, &obj)
	test.NoError(t, err)

	test.Assert(t, obj.Ct == 2)
}

type UID string

type WordCount struct {
	R   api.Request `url:"/wc"`
	UID UID         `api:"auth,required" json:"uid"`
	Str string      `api:"in,required" json:"str"`
	Ct  int         `api:"out" json:"ct"`
}

func (this *WordCount) Do(ctx.C) error {
	this.Ct = len(strings.Fields(this.Str))
	return nil
}
