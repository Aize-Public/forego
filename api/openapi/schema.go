package openapi

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

type Schema struct {
	Type            string             `json:"type,omitempty"`                 // object/string/integer
	Format          string             `json:"format,omitempty"`               // int64
	Required        []string           `json:"required,omitempty"`             // list of required fields names
	Properties      map[string]*Schema `json:"properties,omitempty"`           // for structs
	AdditionalProps *Schema            `json:"additionalProperties,omitempty"` // for maps
	Description     string             `json:"description,omitempty"`
	Example         string             `json:"example,omitempty"`
	Items           *Schema            `json:"items,omitempty"` // for arrays/slices
	Default         string             `json:"default,omitempty"`
	Enum            []string           `json:"enum,omitempty"`

	// Mutually exclusive (if you have a $ref, it will overwrite anything else)
	// See https://swagger.io/docs/specification/using-ref/
	Reference string `json:"$ref,omitempty"`
}

func (this *Service) SchemaFromType(c ctx.C, t reflect.Type, tags *reflect.StructTag) (*Schema, error) {
	s, err := this.schemaFromType(c, t)
	if tags != nil {
		if doc := tags.Get("doc"); doc != "" {
			s.Description = doc
		}
		if example := tags.Get("example"); example != "" {
			s.Example = example
		}
	}
	return s, err
}

func (this *Service) schemaFromType(c ctx.C, t reflect.Type) (s *Schema, err error) {
	defer func() {
		if s.Format == s.Type {
			s.Format = ""
		}
	}()
	zero := reflect.New(t).Elem().Interface()
	switch zero.(type) {
	case json.RawMessage:
		return &Schema{
			Format: "arbitrary JSON",
		}, nil
	case []byte:
		return &Schema{
			Format: "bytes",
		}, nil
	}

	switch zero.(type) {
	case json.Marshaler:
		// Handle corner case where a type can be struct, but marshalled to a primitive json type
		// In openapi schema, we want to describe the marshalled json type
		j, err := json.Marshal(zero)
		if err == nil {
			switch j[0] {
			case '{', 'n', '[': // object like, we dig further
			case '"':
				return &Schema{
					Type:   "string",
					Format: t.String(),
				}, nil
			case 'f', 't':
				return &Schema{
					Type:   "boolean",
					Format: t.String(),
				}, nil
			case '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				// it could be either integer or float here, but it's not possible to determine that in all cases
				return &Schema{
					Type:   "numeric",
					Format: t.String(),
				}, nil
			default:
				panic(fmt.Sprintf("%s => %s", t.String(), j))
			}
		}
	}

	tt := t // tt is the reference type, but t stays the same
	for tt.Kind() == reflect.Pointer {
		tt = tt.Elem()
	}

	switch tt.Kind() {
	case reflect.Bool:
		return &Schema{
			Type:   "boolean",
			Format: t.String(),
		}, nil
	case reflect.Float32, reflect.Float64:
		return &Schema{
			Type:   "numeric",
			Format: t.String(),
		}, nil
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		return &Schema{
			Type:   "integer",
			Format: t.String(),
		}, nil
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return &Schema{
			Type:   "integer",
			Format: t.String(),
		}, nil
	case reflect.String:
		return &Schema{
			Type:   "string",
			Format: t.String(),
		}, nil

	case reflect.Array, reflect.Slice:
		elemSchema, err := this.schemaFromType(c, tt.Elem())
		return &Schema{
			Type:   "array",
			Format: t.String(),
			Items:  elemSchema,
		}, err

	case reflect.Map:
		elemSchema, err := this.schemaFromType(c, tt.Elem())
		return &Schema{
			Type:            "object",
			Format:          t.String(),
			AdditionalProps: elemSchema,
		}, err

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
				s, err := this.schemaFromType(c, f.Type)
				if err != nil {
					return &Schema{}, err
				}
				if doc := f.Tag.Get("doc"); doc != "" {
					s.Description = doc
				}
				if example := f.Tag.Get("example"); example != "" {
					s.Example = example
				}
				name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
				if name == "" {
					name = f.Name
				}
				structSchema.Properties[name] = s
			}
			log.Infof(c, "added referenced struct schema %q: %+v", structKey, structSchema)
		}
		return schema, nil

	default:
		return &Schema{}, ctx.NewErrorf(c, "invalid kind: %v", tt.Kind())
	}
}
