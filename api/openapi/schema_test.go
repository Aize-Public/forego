package openapi_test

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/Aize-Public/forego/api/openapi"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

type Obj struct {
	Map  map[string]bool `json:"map" example:"{\"a\":true, \"b\": false}" doc:"Doc test"`
	List []Sub           `json:"list" example:"[{\"string\": \"a\", \"int\": 42}, {\"string\": \"b\", \"timestamp\": \"2009-11-10T23:00:00Z\"}]"`
}

type Sub struct {
	String  string  `json:"string" example:"s" doc:"test"`
	Boolean bool    `json:"boolean" example:"true"`
	Float   float64 `json:"float" example:"not float"`
	Int     int     `json:"int" example:"-1"`
	Int8    int8    `json:"int8" example:"1"`
	Uint64  uint64  `json:"uint64" example:"not int"`

	Bytes         []byte               `json:"bytes" example:"123"`
	Timestamp     time.Time            `json:"timestamp" example:"2023-11-10T23:00:00Z"`
	Custom        CustomInt            `json:"custom" example:"2"`
	Any           any                  `json:"any" example:"{\"hello\":\"world\"}"`
	CustomEmptyIF CustomEmptyInterface `json:"customEmptyIF" example:"{\"hello\":\"world\"}"`
	SomeInterface io.Reader            `json:"someInterface" example:"whatever"`
	Raw           json.RawMessage      `json:"raw" example:"{\"a\": \"b\", \"c\": [1, 2, 3]}"`

	Anon struct {
		ValueA string `json:"valueA" example:"v"`
		ValueB int    `json:"valueB" example:"0"`
	} `json:"anon" example:"{\"valueA\": \"v\", \"valueB\": 3}"`
	MapAnon map[string]struct {
		ValueA string `json:"valueA" example:"v"`
		ValueB int    `json:"valueB" example:"0"`
	} `json:"mapAnon" example:"{\"a\": {\"valueA\": \"a\", \"valueB\": 3}, \"b\": {\"valueA\": \"b\"}}"`
	ArrAnon []struct {
		ValueA string `json:"valueA" example:"v"`
		ValueB int    `json:"valueB" example:"0"`
	} `json:"arrAnon" example:"[{\"valueA\": \"a\", \"valueB\": 3}, {\"valueA\": \"b\"}]"`

	Loop *Obj `json:"loop"` // tests for infinite recursion
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
	t.Logf("Schema: %s", enc.JSON{Indent: true}.Encode(c, enc.MustMarshal(c, sc)))
	test.EqualsGo(t, "#/components/schemas/github.com_Aize-Public_forego_api_openapi_test_Obj", sc.AllOf[0].Reference)
	t.Logf("Components.Schemas: %s", enc.JSON{Indent: true}.Encode(c, enc.MustMarshal(c, s.Components.Schemas)))

	objSchema := s.Components.Schemas["github.com_Aize-Public_forego_api_openapi_test_Obj"]
	subSchema := s.Components.Schemas["github.com_Aize-Public_forego_api_openapi_test_Sub"]
	test.NotNil(t, objSchema)
	test.NotNil(t, subSchema)

	test.EqualsGo(t, "object", objSchema.Type)
	test.EqualsGo(t, "openapi_test.Obj", objSchema.Format)
	test.NotNil(t, objSchema.Properties["map"])
	test.EqualsGo(t, "object", objSchema.Properties["map"].Type)
	test.EqualsGo(t, "Doc test", objSchema.Properties["map"].Description)
	test.EqualsJSON(c, enc.Map{"a": enc.Bool(true), "b": enc.Bool(false)}, objSchema.Properties["map"].Example)
	test.NotNil(t, objSchema.Properties["map"].AdditionalProps)
	test.EqualsGo(t, "boolean", objSchema.Properties["map"].AdditionalProps.Type)
	test.NotNil(t, objSchema.Properties["list"])
	test.EqualsGo(t, "array", objSchema.Properties["list"].Type)
	test.EqualsJSON(c, enc.List{
		enc.Map{"string": enc.String("a"), "int": enc.Integer(42)},
		enc.Map{"string": enc.String("b"), "timestamp": enc.String("2009-11-10T23:00:00Z")}},
		objSchema.Properties["list"].Example)
	test.NotNil(t, objSchema.Properties["list"].Items)
	test.EqualsGo(t, "#/components/schemas/github.com_Aize-Public_forego_api_openapi_test_Sub", objSchema.Properties["list"].Items.AllOf[0].Reference)

	test.EqualsGo(t, "object", subSchema.Type)
	test.EqualsGo(t, "openapi_test.Sub", subSchema.Format)

	test.NotNil(t, subSchema.Properties["string"])
	test.EqualsGo(t, "string", subSchema.Properties["string"].Type)
	test.EqualsGo(t, "test", subSchema.Properties["string"].Description)
	test.EqualsJSON(c, "s", subSchema.Properties["string"].Example)

	test.NotNil(t, subSchema.Properties["boolean"])
	test.EqualsGo(t, "boolean", subSchema.Properties["boolean"].Type)
	test.EqualsGo(t, "", subSchema.Properties["boolean"].Format)
	test.EqualsJSON(c, true, subSchema.Properties["boolean"].Example)

	test.NotNil(t, subSchema.Properties["float"])
	test.EqualsGo(t, "number", subSchema.Properties["float"].Type)
	test.EqualsGo(t, "float64", subSchema.Properties["float"].Format)
	test.EqualsJSON(c, "not float", subSchema.Properties["float"].Example)
	test.NotNil(t, subSchema.Properties["int"])
	test.EqualsGo(t, "number", subSchema.Properties["int"].Type)
	test.EqualsGo(t, "int", subSchema.Properties["int"].Format)
	test.EqualsJSON(c, "-1", subSchema.Properties["int"].Example)
	test.NotNil(t, subSchema.Properties["int8"])
	test.EqualsGo(t, "number", subSchema.Properties["int8"].Type)
	test.EqualsGo(t, "int8", subSchema.Properties["int8"].Format)
	test.EqualsJSON(c, "1", subSchema.Properties["int8"].Example)
	test.NotNil(t, subSchema.Properties["uint64"])
	test.EqualsGo(t, "number", subSchema.Properties["uint64"].Type)
	test.EqualsGo(t, "uint64", subSchema.Properties["uint64"].Format)
	test.EqualsJSON(c, "not int", subSchema.Properties["uint64"].Example)

	test.NotNil(t, subSchema.Properties["bytes"])
	test.EqualsGo(t, "string", subSchema.Properties["bytes"].Type)
	test.EqualsGo(t, "byte", subSchema.Properties["bytes"].Format)
	test.EqualsJSON(c, "123", subSchema.Properties["bytes"].Example)
	test.NotNil(t, subSchema.Properties["timestamp"])
	test.EqualsGo(t, "string", subSchema.Properties["timestamp"].Type)
	test.EqualsGo(t, "date-time", subSchema.Properties["timestamp"].Format)
	test.EqualsJSON(c, "2023-11-10T23:00:00Z", subSchema.Properties["timestamp"].Example)
	test.NotNil(t, subSchema.Properties["custom"])
	test.EqualsGo(t, "number", subSchema.Properties["custom"].Type)
	test.EqualsGo(t, "openapi_test.CustomInt", subSchema.Properties["custom"].Format)
	test.EqualsJSON(c, 2, subSchema.Properties["custom"].Example)
	test.NotNil(t, subSchema.Properties["any"])
	test.EqualsGo(t, "object", subSchema.Properties["any"].Type)
	test.EqualsGo(t, "", subSchema.Properties["any"].Format)
	test.EqualsJSON(c, enc.Map{"hello": enc.String("world")}, subSchema.Properties["any"].Example)
	test.NotNil(t, subSchema.Properties["customEmptyIF"])
	test.EqualsGo(t, "object", subSchema.Properties["customEmptyIF"].Type)
	test.EqualsGo(t, "openapi_test.CustomEmptyInterface", subSchema.Properties["customEmptyIF"].Format)
	test.EqualsJSON(c, enc.Map{"hello": enc.String("world")}, subSchema.Properties["customEmptyIF"].Example)
	test.NotNil(t, subSchema.Properties["someInterface"])
	test.EqualsGo(t, "object", subSchema.Properties["someInterface"].Type)
	test.EqualsGo(t, "io.Reader", subSchema.Properties["someInterface"].Format)
	test.EqualsJSON(c, "whatever", subSchema.Properties["someInterface"].Example)
	test.NotNil(t, subSchema.Properties["raw"])
	test.EqualsGo(t, "object", subSchema.Properties["raw"].Type)
	test.EqualsGo(t, "", subSchema.Properties["raw"].Format)
	test.EqualsJSON(c, enc.Map{"a": enc.String("b"), "c": enc.List{enc.Integer(1), enc.Integer(2), enc.Integer(3)}}, subSchema.Properties["raw"].Example)

	verifyAnonymous := func(name string, schema *openapi.Schema) {
		t.Logf("Anonymous schema %q", name)
		test.NotNil(t, schema)
		test.EqualsGo(t, "object", schema.Type)
		test.EqualsGo(t, "{valueA string, valueB int}", schema.Format)
		test.NotNil(t, schema.Properties["valueA"])
		test.EqualsGo(t, "string", schema.Properties["valueA"].Type)
		test.EqualsJSON(c, "v", schema.Properties["valueA"].Example)
		test.NotNil(t, schema.Properties["valueB"])
		test.EqualsGo(t, "number", schema.Properties["valueB"].Type)
		test.EqualsJSON(c, 0, schema.Properties["valueB"].Example)
	}
	verifyAnonymous("anon", subSchema.Properties["anon"])
	test.EqualsJSON(c, enc.Map{"valueA": enc.String("v"), "valueB": enc.Integer(3)}, subSchema.Properties["anon"].Example)
	test.NotNil(t, subSchema.Properties["mapAnon"])
	test.EqualsGo(t, "object", subSchema.Properties["mapAnon"].Type)
	test.EqualsGo(t, "map[string]{valueA string, valueB int}", subSchema.Properties["mapAnon"].Format)
	test.EqualsJSON(c, enc.Map{
		"a": enc.Map{"valueA": enc.String("a"), "valueB": enc.Integer(3)},
		"b": enc.Map{"valueA": enc.String("b")}},
		subSchema.Properties["mapAnon"].Example)
	verifyAnonymous("mapAnon", subSchema.Properties["mapAnon"].AdditionalProps)
	test.NotNil(t, subSchema.Properties["arrAnon"])
	test.EqualsGo(t, "array", subSchema.Properties["arrAnon"].Type)
	test.EqualsGo(t, "[]{valueA string, valueB int}", subSchema.Properties["arrAnon"].Format)
	test.EqualsJSON(c, enc.List{
		enc.Map{"valueA": enc.String("a"), "valueB": enc.Integer(3)},
		enc.Map{"valueA": enc.String("b")}},
		subSchema.Properties["arrAnon"].Example)
	verifyAnonymous("arrAnon", subSchema.Properties["arrAnon"].Items)

	test.NotNil(t, s.Components.SecurityScheme)
	test.NotNil(t, s.Components.SecurityScheme["jwt"])
	test.EqualsGo(t, "http", s.Components.SecurityScheme["jwt"].Type)
	test.EqualsGo(t, "bearer", s.Components.SecurityScheme["jwt"].Scheme)
	test.EqualsGo(t, "JWT", s.Components.SecurityScheme["jwt"].BearerFormat)
}
