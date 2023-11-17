package openapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
)

type Schema struct {
	Type            string             `json:"type,omitempty"`
	Format          string             `json:"format,omitempty"`
	Required        []string           `json:"required,omitempty"`             // list of required fields names
	Properties      map[string]*Schema `json:"properties,omitempty"`           // for structs
	AdditionalProps *Schema            `json:"additionalProperties,omitempty"` // for maps
	Description     string             `json:"description,omitempty"`
	Example         any                `json:"example,omitempty"`
	Items           *Schema            `json:"items,omitempty"` // for arrays/slices
	Default         string             `json:"default,omitempty"`
	Enum            []string           `json:"enum,omitempty"`
	AllOf           []*Schema          `json:"allOf,omitempty"`

	// Mutually exclusive (if you have a $ref, it will overwrite anything else)
	// See https://swagger.io/docs/specification/using-ref/
	Reference string `json:"$ref,omitempty"`
}

func (this *Service) SchemaFromType(c ctx.C, t reflect.Type, tags *reflect.StructTag) (*Schema, error) {
	var doc, example string
	if tags != nil {
		doc = tags.Get("doc")
		example = strings.TrimSpace(tags.Get("example"))
	}
	s, err := this.schemaFromType(c, t, doc, example)
	return s, err
}

func (this *Service) schemaFromType(c ctx.C, t reflect.Type, doc, example string) (s *Schema, err error) {
	defer func() {
		if s.Format == s.Type {
			s.Format = ""
		}
		s.Description = doc
		if example == "" {
			s.Example = nil
		}
	}()

	tt := t                            // tt is the reference type, but t stays the same
	for tt.Kind() == reflect.Pointer { // if pointer, find its value type
		tt = tt.Elem()
	}

	zero := reflect.New(tt).Elem().Interface()
	switch zero.(type) {
	case json.RawMessage:
		return &Schema{
			Type:    "object",
			Format:  "",
			Example: example,
		}, nil

	case []byte:
		return &Schema{
			Type:    "string",
			Format:  "byte",
			Example: example,
		}, nil

	case time.Time, enc.Time:
		return &Schema{
			Type:    "string",
			Format:  "date-time",
			Example: example,
		}, nil

	case json.Marshaler, enc.Marshaler:
		// Handle corner case where a type can be struct, but marshalled to a primitive json type
		// In openapi schema, we want to describe the marshalled json type
		j, err := enc.MarshalJSON(c, zero)
		if err == nil {
			switch j[0] {
			case '{', 'n', '[':
				// object like, we dig further...
			case '"':
				return &Schema{
					Type:    "string",
					Format:  t.String(),
					Example: example,
				}, nil
			case 'f', 't':
				return &Schema{
					Type:    "boolean",
					Format:  t.String(),
					Example: tryDecodingAsBool(c, example),
				}, nil
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
				// It could be either integer or float here, but it's not always possible to determine which
				// Therefore just setting the type to "number", and doing best effort parsing of the example
				var ex any
				if strings.Contains(example, ".") {
					ex = tryDecodingAsFloat(c, example)
				} else {
					ex = tryDecodingAsInt(c, example)
				}
				return &Schema{
					Type:    "number",
					Format:  t.String(),
					Example: ex,
				}, nil
			default:
				panic(fmt.Sprintf("%s => %s", t.String(), j))
			}
		}
	}

	switch tt.Kind() {
	case reflect.Bool:
		return &Schema{
			Type:    "boolean",
			Format:  "",
			Example: tryDecodingAsBool(c, example),
		}, nil

	case reflect.Float32, reflect.Float64:
		return &Schema{
			Type:    "number",
			Format:  t.String(),
			Example: tryDecodingAsFloat(c, example),
		}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fallthrough
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &Schema{
			Type:    "number", // Using "number" instead of "integer" here as well, to be consistent
			Format:  t.String(),
			Example: tryDecodingAsInt(c, example),
		}, nil

	case reflect.String:
		return &Schema{
			Type:    "string",
			Format:  t.String(),
			Example: example,
		}, nil

	case reflect.Array, reflect.Slice:
		elemSchema, err := this.schemaFromType(c, tt.Elem(), "", "")
		return &Schema{
			Type:    "array",
			Format:  t.String(),
			Items:   elemSchema,
			Example: tryDecodingAsList(c, example),
		}, err

	case reflect.Map:
		elemSchema, err := this.schemaFromType(c, tt.Elem(), "", "")
		return &Schema{
			Type:            "object",
			Format:          t.String(),
			AdditionalProps: elemSchema,
			Example:         tryDecodingAsMap(c, example),
		}, err

	case reflect.Interface:
		format := ""
		if t != reflect.TypeOf((*any)(nil)).Elem() {
			// Just exposing the name of custom interfaces, but maybe we could do something better here
			format = t.String()
		}
		return &Schema{
			Type:    "object",
			Format:  format,
			Example: tryDecodingAsNode(c, example),
		}, nil

	case reflect.Struct:
		// Create reference
		structKey := tt.PkgPath() + "/" + tt.Name()
		structKey = strings.ReplaceAll(structKey, "/", "_")
		schema := &Schema{Reference: "#/components/schemas/" + structKey}

		// Add definition if not exists (the same struct type should only be defined once)
		if this.Components.Schemas == nil {
			this.Components.Schemas = make(map[string]*Schema)
		}
		if _, exists := this.Components.Schemas[structKey]; !exists {
			structSchema := &Schema{
				Type:       "object",
				Format:     t.String(),
				Properties: map[string]*Schema{},
			}
			this.Components.Schemas[structKey] = structSchema // adding it before recursion, to protect against infinite recursion

			for i := 0; i < tt.NumField(); i++ {
				f := tt.Field(i)
				s, err := this.schemaFromType(c, f.Type, f.Tag.Get("doc"), strings.TrimSpace(f.Tag.Get("example")))
				if err != nil {
					return &Schema{}, err
				}
				name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
				if name == "" {
					name = f.Name
				}
				structSchema.Properties[name] = s
			}
			log.Infof(c, "added referenced struct schema %q: %+v", structKey, structSchema)
		}

		// Using "allOf" as a workaround for not being able to add example/description to a $ref
		return &Schema{
			Example: tryDecodingAsMap(c, example),
			AllOf:   []*Schema{schema},
		}, nil

	default:
		return &Schema{}, ctx.NewErrorf(c, "invalid kind: %v", tt.Kind())
	}
}

func tryDecodingAsBool(c ctx.C, v string) any {
	switch v {
	case "":
		return nil
	case "true":
		return true
	case "false":
		return false
	}
	log.Warnf(c, "Unable to parse example %q as boolean", v)
	return v
}

func tryDecodingAsFloat(c ctx.C, v string) any {
	if v == "" {
		return nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		log.Warnf(c, "Unable to parse example %q as float: %s", v, err)
		return v
	}
	return f
}

func tryDecodingAsInt(c ctx.C, v string) any {
	if v == "" {
		return nil
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		log.Warnf(c, "Unable to parse example %q as int: %s", v, err)
		return v
	}
	return int(i)
}

func tryDecodingAsList(c ctx.C, v string) any {
	if v == "" {
		return nil
	}
	if node, err := (enc.JSON{}).Decode(c, []byte(v)); err == nil {
		if l, ok := node.(enc.List); ok {
			return l
		}
	}
	return v
}

func tryDecodingAsMap(c ctx.C, v string) any {
	if v == "" {
		return nil
	}
	if node, err := (enc.JSON{}).Decode(c, []byte(v)); err == nil {
		if m, ok := node.(enc.Map); ok {
			return m
		}
	}
	return v
}

func tryDecodingAsNode(c ctx.C, v string) any {
	if v == "" {
		return nil
	}
	if node, err := (enc.JSON{}).Decode(c, []byte(v)); err == nil {
		return node
	}
	return v
}
