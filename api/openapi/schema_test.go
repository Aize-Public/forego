package openapi_test

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/Aize-Public/forego/api/openapi"
	"github.com/Aize-Public/forego/test"
)

type Obj struct {
	Map  map[string]bool `json:"map"`
	List []Sub           `json:"list"`
}

type Sub struct {
	String  string  `json:"string"`
	Boolean bool    `json:"boolean"`
	Float   float64 `json:"float"`
	Int     int     `json:"int"`
	Int8    int8    `json:"int8"`
	Uint64  uint64  `json:"uint64"`

	Bytes         []byte               `json:"bytes"`
	Timestamp     time.Time            `json:"timestamp"`
	Custom        CustomInt            `json:"custom"`
	Any           any                  `json:"any"`
	CustomEmptyIF CustomEmptyInterface `json:"customEmptyIF"`
	SomeInterface io.Reader            `json:"someInterface"`
	Raw           json.RawMessage      `json:"raw"`
}

type CustomInt struct {
	Value int
}

func (this CustomInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.Value)
}

type CustomEmptyInterface interface{}

func TestSchema(t *testing.T) {
	c := test.Context(t)

	s := openapi.NewService("test-schema")
	sc, err := s.SchemaFromType(c, reflect.TypeOf(Obj{}), nil)
	test.NoError(t, err)
	t.Logf("Schema: %+v", sc)
	test.EqualsGo(t, "#/components/schemas/github.com_Aize-Public_forego_api_openapi_test_Obj", sc.Reference)
	t.Logf("Components.Schemas: %+v", s.Components.Schemas)

	objSchema := s.Components.Schemas["github.com_Aize-Public_forego_api_openapi_test_Obj"]
	subSchema := s.Components.Schemas["github.com_Aize-Public_forego_api_openapi_test_Sub"]
	test.NotNil(t, objSchema)
	test.NotNil(t, subSchema)

	test.EqualsGo(t, "object", objSchema.Type)
	test.EqualsGo(t, "openapi_test.Obj", objSchema.Format)
	test.NotNil(t, objSchema.Properties["map"])
	test.EqualsGo(t, "object", objSchema.Properties["map"].Type)
	test.NotNil(t, objSchema.Properties["map"].AdditionalProps)
	test.EqualsGo(t, "boolean", objSchema.Properties["map"].AdditionalProps.Type)
	test.NotNil(t, objSchema.Properties["list"])
	test.EqualsGo(t, "array", objSchema.Properties["list"].Type)
	test.NotNil(t, objSchema.Properties["list"].Items)
	test.EqualsGo(t, "#/components/schemas/github.com_Aize-Public_forego_api_openapi_test_Sub", objSchema.Properties["list"].Items.Reference)

	test.EqualsGo(t, "object", subSchema.Type)
	test.EqualsGo(t, "openapi_test.Sub", subSchema.Format)

	test.NotNil(t, subSchema.Properties["string"])
	test.EqualsGo(t, "string", subSchema.Properties["string"].Type)

	test.NotNil(t, subSchema.Properties["boolean"])
	test.EqualsGo(t, "boolean", subSchema.Properties["boolean"].Type)
	test.EqualsGo(t, "", subSchema.Properties["boolean"].Format)

	test.NotNil(t, subSchema.Properties["float"])
	test.EqualsGo(t, "number", subSchema.Properties["float"].Type)
	test.EqualsGo(t, "float64", subSchema.Properties["float"].Format)
	test.NotNil(t, subSchema.Properties["int"])
	test.EqualsGo(t, "number", subSchema.Properties["int"].Type)
	test.EqualsGo(t, "int", subSchema.Properties["int"].Format)
	test.NotNil(t, subSchema.Properties["int8"])
	test.EqualsGo(t, "number", subSchema.Properties["int8"].Type)
	test.EqualsGo(t, "int8", subSchema.Properties["int8"].Format)
	test.NotNil(t, subSchema.Properties["uint64"])
	test.EqualsGo(t, "number", subSchema.Properties["uint64"].Type)
	test.EqualsGo(t, "uint64", subSchema.Properties["uint64"].Format)

	test.NotNil(t, subSchema.Properties["bytes"])
	test.EqualsGo(t, "string", subSchema.Properties["bytes"].Type)
	test.EqualsGo(t, "byte", subSchema.Properties["bytes"].Format)
	test.NotNil(t, subSchema.Properties["timestamp"])
	test.EqualsGo(t, "string", subSchema.Properties["timestamp"].Type)
	test.EqualsGo(t, "date-time", subSchema.Properties["timestamp"].Format)
	test.NotNil(t, subSchema.Properties["custom"])
	test.EqualsGo(t, "number", subSchema.Properties["custom"].Type)
	test.EqualsGo(t, "openapi_test.CustomInt", subSchema.Properties["custom"].Format)
	test.NotNil(t, subSchema.Properties["any"])
	test.EqualsGo(t, "object", subSchema.Properties["any"].Type)
	test.EqualsGo(t, "", subSchema.Properties["any"].Format)
	test.NotNil(t, subSchema.Properties["customEmptyIF"])
	test.EqualsGo(t, "object", subSchema.Properties["customEmptyIF"].Type)
	test.EqualsGo(t, "openapi_test.CustomEmptyInterface", subSchema.Properties["customEmptyIF"].Format)
	test.NotNil(t, subSchema.Properties["someInterface"])
	test.EqualsGo(t, "object", subSchema.Properties["someInterface"].Type)
	test.EqualsGo(t, "io.Reader", subSchema.Properties["someInterface"].Format)
	test.NotNil(t, subSchema.Properties["raw"])
	test.EqualsGo(t, "object", subSchema.Properties["raw"].Type)
	test.EqualsGo(t, "", subSchema.Properties["raw"].Format)
}
