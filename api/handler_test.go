package api_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/Aize-Public/forego/api"
	"github.com/Aize-Public/forego/ctx"
)

func TestHandler(t *testing.T) {
	c := context.Background()
	type T struct {
		Str string `api:"str,in"`
		Int int    `api:"int,out"`
	}
	h, err := api.NewHandler(c, T{})
	if err != nil {
		t.Fatal(err)
	}
	data := TestData{
		AuthFunc: func(c ctx.C, into reflect.Value, required bool) error {
			return nil
		},
		UnmarshalFunc: func(c ctx.C, name string, into reflect.Value) error {
			switch name {
			case "str":
				into.SetString("foo")
			default:
			}
			t.Logf("unmarshal(%q) => %v", name, into)
			return nil
		},
		MarshalFunc: func(c ctx.C, name string, from reflect.Value) error {
			t.Logf("marshal(%q) <= %v", name, from)
			return nil
		},
	}
	obj, err := h.Server().Recv(c, data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("obj: %+v", obj)
	obj.Int = len(obj.Str)
	err = h.Server().Send(c, obj, data)
	if err != nil {
		t.Fatal(err)
	}
}

type TestData struct {
	AuthFunc      func(c ctx.C, into reflect.Value, required bool) error
	UnmarshalFunc func(c ctx.C, name string, into reflect.Value) error
	MarshalFunc   func(c ctx.C, name string, from reflect.Value) error
}

func (this TestData) Auth(c ctx.C, into reflect.Value, required bool) error {
	return this.AuthFunc(c, into, required)
}

func (this TestData) Unmarshal(c ctx.C, name string, into reflect.Value) error {
	return this.UnmarshalFunc(c, name, into)
}

func (this TestData) Marshal(c ctx.C, name string, into reflect.Value) error {
	return this.MarshalFunc(c, name, into)
}

var _ api.ClientRequest = TestData{}
