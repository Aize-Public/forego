package openapi_test

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/Aize-Public/forego/api/openapi"
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/enc"
	"github.com/Aize-Public/forego/test"
)

// Test object
type Obj struct {
	SimpleMap map[string]bool `json:"simpleMap" example:"{\"a\":true, \"b\": false}" doc:"Doc test"`
	Map       map[string]Sub  `json:"map" example:"{\"a\": {\"string\": \"a\", \"int\": 42}}"`
	List      []Sub           `json:"list" example:"[{\"string\": \"a\", \"int\": 42}, {\"string\": \"b\", \"timestamp\": \"2009-11-10T23:00:00Z\"}]"`
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
	EncMap        enc.Map              `json:"encMap" example:"{\"a\":true, \"b\": 42, \"c\": \"s\"}"`
	EncList       enc.List             `json:"encList" example:"[true, 42, \"s\"]"`

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
	test.NotNil(t, objSchema)
	test.EqualsGo(t, "object", objSchema.Type)
	test.EqualsGo(t, "openapi_test.Obj", objSchema.Format)

	{
		s := objSchema.Properties["simpleMap"]
		test.NotNil(t, s)
		test.EqualsGo(t, "object", s.Type)
		test.EqualsGo(t, "map[string]boolean", s.Format)
		test.EqualsGo(t, "Doc test", s.Description)
		test.EqualsJSON(c, enc.Map{"a": enc.Bool(true), "b": enc.Bool(false)}, s.Example)
		test.NotNil(t, s.AdditionalProps)
		test.EqualsGo(t, "boolean", s.AdditionalProps.Type)
	}
	{
		s := objSchema.Properties["map"]
		test.NotNil(t, s)
		test.EqualsGo(t, "object", s.Type)
		test.EqualsGo(t, "map[string]openapi_test.Sub", s.Format)
		test.EqualsJSON(c, enc.Map{"a": enc.Map{"string": enc.String("a"), "int": enc.Integer(42)}}, s.Example)
		test.NotNil(t, s.AdditionalProps)
		test.EqualsGo(t, "#/components/schemas/github.com_Aize-Public_forego_api_openapi_test_Sub", s.AdditionalProps.AllOf[0].Reference)
	}
	{
		s := objSchema.Properties["list"]
		test.NotNil(t, s)
		test.EqualsGo(t, "array", s.Type)
		test.EqualsGo(t, "[]openapi_test.Sub", s.Format)
		test.EqualsJSON(c, enc.List{
			enc.Map{"string": enc.String("a"), "int": enc.Integer(42)},
			enc.Map{"string": enc.String("b"), "timestamp": enc.String("2009-11-10T23:00:00Z")}},
			s.Example)
		test.NotNil(t, s.Items)
		test.EqualsGo(t, "#/components/schemas/github.com_Aize-Public_forego_api_openapi_test_Sub", s.Items.AllOf[0].Reference)
	}
	{
		s := s.Components.Schemas["github.com_Aize-Public_forego_api_openapi_test_Sub"]
		verifySubSchema(c, t, s)
	}

	test.NotNil(t, s.Components.SecurityScheme)
	test.NotNil(t, s.Components.SecurityScheme["jwt"])
	test.EqualsGo(t, "http", s.Components.SecurityScheme["jwt"].Type)
	test.EqualsGo(t, "bearer", s.Components.SecurityScheme["jwt"].Scheme)
	test.EqualsGo(t, "JWT", s.Components.SecurityScheme["jwt"].BearerFormat)
}

func verifySubSchema(c ctx.C, t *testing.T, schema *openapi.Schema) {
	test.NotNil(t, schema)
	test.EqualsGo(t, "object", schema.Type)
	test.EqualsGo(t, "openapi_test.Sub", schema.Format)

	props := schema.Properties

	test.NotNil(t, props["string"])
	test.EqualsGo(t, "string", props["string"].Type)
	test.EqualsGo(t, "", props["string"].Format)
	test.EqualsGo(t, "test", props["string"].Description)
	test.EqualsJSON(c, "s", props["string"].Example)

	test.NotNil(t, props["boolean"])
	test.EqualsGo(t, "boolean", props["boolean"].Type)
	test.EqualsGo(t, "", props["boolean"].Format)
	test.EqualsJSON(c, true, props["boolean"].Example)

	test.NotNil(t, props["float"])
	test.EqualsGo(t, "number", props["float"].Type)
	test.EqualsGo(t, "float64", props["float"].Format)
	test.EqualsJSON(c, "not float", props["float"].Example)
	test.NotNil(t, props["int"])
	test.EqualsGo(t, "number", props["int"].Type)
	test.EqualsGo(t, "int", props["int"].Format)
	test.EqualsJSON(c, "-1", props["int"].Example)
	test.NotNil(t, props["int8"])
	test.EqualsGo(t, "number", props["int8"].Type)
	test.EqualsGo(t, "int8", props["int8"].Format)
	test.EqualsJSON(c, "1", props["int8"].Example)
	test.NotNil(t, props["uint64"])
	test.EqualsGo(t, "number", props["uint64"].Type)
	test.EqualsGo(t, "uint64", props["uint64"].Format)
	test.EqualsJSON(c, "not int", props["uint64"].Example)

	test.NotNil(t, props["bytes"])
	test.EqualsGo(t, "string", props["bytes"].Type)
	test.EqualsGo(t, "byte", props["bytes"].Format)
	test.EqualsJSON(c, "123", props["bytes"].Example)
	test.NotNil(t, props["timestamp"])
	test.EqualsGo(t, "string", props["timestamp"].Type)
	test.EqualsGo(t, "date-time", props["timestamp"].Format)
	test.EqualsJSON(c, "2023-11-10T23:00:00Z", props["timestamp"].Example)
	test.NotNil(t, props["custom"])
	test.EqualsGo(t, "number", props["custom"].Type)
	test.EqualsGo(t, "openapi_test.CustomInt", props["custom"].Format)
	test.EqualsJSON(c, 2, props["custom"].Example)
	test.NotNil(t, props["any"])
	test.EqualsGo(t, "object", props["any"].Type)
	test.EqualsGo(t, "", props["any"].Format)
	test.EqualsJSON(c, enc.Map{"hello": enc.String("world")}, props["any"].Example)
	test.NotNil(t, props["customEmptyIF"])
	test.EqualsGo(t, "object", props["customEmptyIF"].Type)
	test.EqualsGo(t, "openapi_test.CustomEmptyInterface", props["customEmptyIF"].Format)
	test.EqualsJSON(c, enc.Map{"hello": enc.String("world")}, props["customEmptyIF"].Example)
	test.NotNil(t, props["someInterface"])
	test.EqualsGo(t, "object", props["someInterface"].Type)
	test.EqualsGo(t, "io.Reader", props["someInterface"].Format)
	test.EqualsJSON(c, "whatever", props["someInterface"].Example)
	test.NotNil(t, props["raw"])
	test.EqualsGo(t, "object", props["raw"].Type)
	test.EqualsGo(t, "", props["raw"].Format)
	test.EqualsJSON(c, enc.Map{"a": enc.String("b"), "c": enc.List{enc.Integer(1), enc.Integer(2), enc.Integer(3)}}, props["raw"].Example)

	test.NotNil(t, props["encMap"])
	test.EqualsGo(t, "object", props["encMap"].Type)
	test.EqualsGo(t, "map[string]enc.Node", props["encMap"].Format)
	test.EqualsJSON(c, enc.Map{
		"a": enc.Bool(true),
		"b": enc.Integer(42),
		"c": enc.String("s")},
		props["encMap"].Example)
	test.NotNil(t, props["encList"])
	test.EqualsGo(t, "array", props["encList"].Type)
	test.EqualsGo(t, "[]enc.Node", props["encList"].Format)
	test.EqualsJSON(c, enc.List{
		enc.Bool(true),
		enc.Integer(42),
		enc.String("s")},
		props["encList"].Example)

	verifyAnonymous := func(name string, s *openapi.Schema) {
		t.Logf("Anonymous schema %q", name)
		test.NotNil(t, s)
		test.EqualsGo(t, "object", s.Type)
		test.EqualsGo(t, "{valueA string, valueB int}", s.Format)
		test.NotNil(t, s.Properties["valueA"])
		test.EqualsGo(t, "string", s.Properties["valueA"].Type)
		test.EqualsJSON(c, "v", s.Properties["valueA"].Example)
		test.NotNil(t, s.Properties["valueB"])
		test.EqualsGo(t, "number", s.Properties["valueB"].Type)
		test.EqualsJSON(c, 0, s.Properties["valueB"].Example)
	}
	verifyAnonymous("anon", props["anon"])
	test.EqualsJSON(c, enc.Map{"valueA": enc.String("v"), "valueB": enc.Integer(3)}, props["anon"].Example)
	test.NotNil(t, props["mapAnon"])
	test.EqualsGo(t, "object", props["mapAnon"].Type)
	test.EqualsGo(t, "map[string]{valueA string, valueB int}", props["mapAnon"].Format)
	test.EqualsJSON(c, enc.Map{
		"a": enc.Map{"valueA": enc.String("a"), "valueB": enc.Integer(3)},
		"b": enc.Map{"valueA": enc.String("b")}},
		props["mapAnon"].Example)
	verifyAnonymous("mapAnon", props["mapAnon"].AdditionalProps)
	test.NotNil(t, props["arrAnon"])
	test.EqualsGo(t, "array", props["arrAnon"].Type)
	test.EqualsGo(t, "[]{valueA string, valueB int}", props["arrAnon"].Format)
	test.EqualsJSON(c, enc.List{
		enc.Map{"valueA": enc.String("a"), "valueB": enc.Integer(3)},
		enc.Map{"valueA": enc.String("b")}},
		props["arrAnon"].Example)
	verifyAnonymous("arrAnon", props["arrAnon"].Items)
}
