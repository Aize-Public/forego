package enc

import (
	"fmt"
	"reflect"
	"strings"
)

type Tag struct {
	Name      string
	OmitEmpty bool
	Skip      bool
}

// TODO(oha): we need to parse `enc`, `json` and eventually `yaml` and make sure the agree
func parseTag(tag reflect.StructField) (out Tag) {
	json := tag.Tag.Get("json")
	if json == "-" {
		out.Skip = true
		return
	}
	name, extra, _ := strings.Cut(json, ",")
	if name == "" {
		name = tag.Name
	}
	out.Name = name
	switch extra {
	case "omitempty":
		out.OmitEmpty = true
	case "":
	default:
		panic(fmt.Sprintf("invalid tag: %v", tag))
	}
	return
}
