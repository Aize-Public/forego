package api_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/Aize-Public/forego/api"
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
		UnmarshalFunc: func(c context.Context, name string, into reflect.Value) error {
			switch name {
			case "str":
				into.SetString("foo")
			default:
			}
			t.Logf("unmarshal(%q) => %v", name, into)
			return nil
		},
		MarshalFunc: func(c context.Context, name string, from reflect.Value) error {
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
	UnmarshalFunc func(c context.Context, name string, into reflect.Value) error
	MarshalFunc   func(c context.Context, name string, from reflect.Value) error
}

func (this TestData) Unmarshal(c context.Context, name string, into reflect.Value) error {
	return this.UnmarshalFunc(c, name, into)
}

func (this TestData) Marshal(c context.Context, name string, into reflect.Value) error {
	return this.MarshalFunc(c, name, into)
}

var _ api.Data = TestData{}
