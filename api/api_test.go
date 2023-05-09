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
		UnmarshalFunc: func(c ctx.C, name string, into any) error {
			switch name {
			case "str":
				reflect.ValueOf(into).SetString("foo")
			default:
			}
			t.Logf("unmarshal(%q) => %v", name, into)
			return nil
		},
		MarshalFunc: func(c ctx.C, name string, from any) error {
			t.Logf("marshal(%q) <= %v", name, from)
			return nil
		},
	}
	obj, err := h.RequestIn(c, data)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("obj: %+v", obj)
	obj.Int = len(obj.Str)
	err = h.ResponseOut(c, obj, data)
	if err != nil {
		t.Fatal(err)
	}
}

type TestData struct {
	UnmarshalFunc func(c ctx.C, name string, into any) error
	MarshalFunc   func(c ctx.C, name string, from any) error
}

func (this TestData) Unmarshal(c ctx.C, name string, into any) error {
	return this.UnmarshalFunc(c, name, into)
}

func (this TestData) Marshal(c ctx.C, name string, into any) error {
	return this.MarshalFunc(c, name, into)
}

var _ api.Data = TestData{}
