package ctx

import (
	"encoding/json"
	"fmt"
)

type tagC struct {
	C
	key  string
	json JSON
}

// Append the given key/val pair to the context
// Note(oha): if the same key is added twice, it will be found twice later with the Range
func WithTag(c C, key string, val any) C {
	var j []byte
	switch val := val.(type) {
	case json.RawMessage:
		j = val
	case JSON:
		j = val
	case []byte:
		if json.Valid(val) {
			j = val
		}
	default:
	}
	if j == nil {
		var err error
		j, err = json.Marshal(val)
		if err != nil {
			j, _ = json.Marshal(fmt.Sprintf("can't tag type %T: %v", val, err))
		}
	}
	return tagC{
		C:    c,
		key:  key,
		json: j,
	}
}

// scan the chain of context for tags, call the function on each of them, parents first
func RangeTag(c C, fn func(k string, json JSON) error) error {
	v := c.Value(tagRangeFunc(fn))
	switch v := v.(type) {
	case nil:
		return nil
	case error:
		return v
	default:
		return fmt.Errorf("error-ish: %+v", v)
	}
}

type tagRangeFunc func(k string, json JSON) error

func (c tagC) Value(k any) any {
	switch k := k.(type) {
	case tagRangeFunc:
		//log.Printf("tagC.Value... %p", k)
		err := c.C.Value(k) // parents first
		if err != nil {
			return err
		}
		return k(c.key, c.json)
	default:
		//log.Printf("tagC.Value(%T)...", k)
		return c.C.Value(k)
	}
}
