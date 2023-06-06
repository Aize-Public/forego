package enc

import (
	"fmt"
	"reflect"
	"strings"
)

type Tag struct {
	Name      string
	JSON      string
	OmitEmpty bool
	Skip      bool
}

// TODO(oha): we need to parse `enc`, `json` and eventually `yaml` and make sure they agree
func parseTag(tag reflect.StructField) (out Tag) {
	out.Name = tag.Name

	json := tag.Tag.Get("json")
	if json == "-" {
		out.Skip = true
		return
	}
	json, extra, _ := strings.Cut(json, ",")
	if json == "" {
		json = tag.Name // NOTE(oha): should we use "enc" name?
	}
	out.JSON = json
	switch extra {
	case "omitempty":
		out.OmitEmpty = true
	case "":
	default:
		panic(fmt.Sprintf("invalid tag: %v", tag))
	}
	return
}
