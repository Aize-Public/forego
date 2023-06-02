package api_test

import (
	"strings"
	"testing"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

func TestAPI(t *testing.T) {
	c := test.Context(t)
	alice := UID("alice")

	t.Logf("handler...")
	ser, err := api.NewServer(c, &WordCount{})
	test.NoError(t, err)

	cli, err := api.NewClient(c, &WordCount{})
	test.NoError(t, err)

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
	data.UID, _ = enc.Marshal(c, alice)
	//data.UID, _ = json.Marshal(alice)

	{
		t.Logf("server recv...")
		obj, err := ser.Recv(c, data)
		test.NoError(t, err)
		test.NotEmpty(t, obj)
		test.EqualsJSON(t, alice, obj.UID)

		t.Logf("Foo()ing...")
		err = obj.Foo(c)
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
	Str string      `api:"in" json:"str"`
	Ct  int         `api:"out" json:"ct"`
}

func (this *WordCount) Foo(ctx.C) error {
	this.Ct = len(strings.Split(this.Str, " "))
	return nil
}
