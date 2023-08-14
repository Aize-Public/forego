package openapi_test

import (
	"reflect"
	"testing"

	"github.com/Aize-Public/forego/api/openapi"
	"github.com/Aize-Public/forego/test"
)

type Obj struct {
	ID  string   `json:"id"`
	Out []Record `json:"out"`
}

type Record struct {
	Score float64 `json:"score"`
	Label string  `json:"label"`
	Blob  any     `json:"blob"`
}

func TestSchema(t *testing.T) {
	c := test.Context(t)

	s := openapi.NewService("test-schema")
	sc, err := s.SchemaFromType(c, reflect.TypeOf(Obj{}), nil)
	test.NoError(t, err)
	t.Logf("Schema: %+v", sc)
	test.EqualsGo(t, "#/components/schemas/github.com_Aize-Public_forego_api_openapi_test_Obj", sc.Reference)
	t.Logf("Components.Schemas: %+v", s.Components.Schemas)

	objSchema := s.Components.Schemas["github.com_Aize-Public_forego_api_openapi_test_Obj"]
	recordSchema := s.Components.Schemas["github.com_Aize-Public_forego_api_openapi_test_Record"]
	test.NotNil(t, objSchema)
	test.NotNil(t, recordSchema)

	test.EqualsGo(t, "object", objSchema.Type)
	test.EqualsGo(t, "openapi_test.Obj", objSchema.Format)
	test.NotNil(t, objSchema.Properties["id"])
	test.EqualsGo(t, "string", objSchema.Properties["id"].Type)
	test.NotNil(t, objSchema.Properties["out"])
	test.EqualsGo(t, "array", objSchema.Properties["out"].Type)
	test.NotNil(t, objSchema.Properties["out"].Items)
	test.EqualsGo(t, "#/components/schemas/github.com_Aize-Public_forego_api_openapi_test_Record", objSchema.Properties["out"].Items.Reference)

	test.EqualsGo(t, "object", recordSchema.Type)
	test.EqualsGo(t, "openapi_test.Record", recordSchema.Format)
	test.NotNil(t, recordSchema.Properties["score"])
	test.EqualsGo(t, "numeric", recordSchema.Properties["score"].Type)
	test.EqualsGo(t, "float64", recordSchema.Properties["score"].Format)
	test.NotNil(t, recordSchema.Properties["label"])
	test.EqualsGo(t, "string", recordSchema.Properties["label"].Type)
	test.NotNil(t, recordSchema.Properties["blob"])
	test.EqualsGo(t, "object", recordSchema.Properties["blob"].Type)
	test.EqualsGo(t, "interface {}", recordSchema.Properties["blob"].Format)
}
