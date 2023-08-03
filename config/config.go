package config

import (
	"encoding"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"

	"github.com/Aize-Public/forego/ctx"
)

func Must[T any](c ctx.C, cfg T, f func(string) string) T {
	cfg, err := From(c, cfg, f)
	if err != nil {
		panic(err)
	}
	return cfg
}

func From[T any](c ctx.C, cfg T, f func(string) string) (T, error) {
	if reflect.TypeOf(cfg).Kind() != reflect.Struct {
		return cfg, ctx.NewErrorf(c, "expected struct, got %T", cfg)
	}
	v := reflect.ValueOf(&cfg).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		ft := t.Field(i)
		tag := ft.Tag.Get("config")
		if tag == "" {
			continue
		}
		parts := strings.Split(tag, ",")
		key := parts[0]

		var defVal string
		var hasDef bool
		for _, def := range parts[1:] {
			parts := strings.SplitN(def, "=", 2)
			if len(parts) == 1 {
				parts = append(parts, "")
			}
			switch parts[0] {
			case "default":
				defVal = parts[1]
				hasDef = true
			default:
				return cfg, ctx.NewErrorf(c, "unsupported tag definition for %q: %q", key, def)
			}
		}
		val := f(key)
		if val == "" {
			if !hasDef {
				return cfg, ctx.NewErrorf(c, "missing config for %q", key)
			}
			val = defVal
		}
		fv := v.Field(i)
		err := unmarshalText(c, fv, key, val)
		if err != nil {
			return cfg, err
		}
		//log.Debugf(c, "config: %q=%#v", key, fv)
	}
	return cfg, nil
}

type UnmarshalText interface {
	UnmarshalText(c ctx.C, t string) error
}

func unmarshalText(c ctx.C, dest reflect.Value, name string, val string) (err error) {
	v := dest.Addr().Interface()
	switch v := v.(type) {
	case UnmarshalText:
		return v.UnmarshalText(c, val)
	case encoding.TextUnmarshaler:
		if val != "" {
			return ctx.WrapError(c, v.UnmarshalText([]byte(val)))
		}
	case json.Unmarshaler:
		if val != "" {
			return ctx.WrapError(c, v.UnmarshalJSON([]byte(val)))
		}
	default:
		//fmt.Printf("OHA %T\n", v)
	}
	switch dest.Kind() {
	case reflect.Bool:
		var v bool
		switch strings.ToLower(val) {
		case "true", "1", "yes", "y":
			v = true
		case "", "false", "0", "no", "n":
			v = false
		default:
			return ctx.NewErrorf(nil, "can't convert to boobool field %s: %q", name, val)
		}
		dest.Set(reflect.ValueOf(v))
	case reflect.Int:
		i, err := strconv.Atoi(val)
		if err != nil {
			return ctx.NewErrorf(nil, "can't convert to int field %s: %q", name, val)
		}
		dest.Set(reflect.ValueOf(i))
	case reflect.String:
		v := reflect.ValueOf(val).Convert(dest.Type())
		dest.Set(v)
	case reflect.Slice:
		if val != "" {
			parts := strings.Split(val, ",")
			x := dest
			for _, el := range parts {
				v := reflect.Indirect(reflect.New(dest.Type().Elem()))
				err := unmarshalText(c, v, name+"[]", el)
				if err != nil {
					return err
				}
				x = reflect.Append(x, v)
			}
			dest.Set(x)
		}

	case reflect.Map:
		if val != "" {
			m := reflect.MakeMap(dest.Type())
			parts := strings.Split(val, ",")
			for _, el := range parts {
				key, val, _ := strings.Cut(el, "=")
				v := reflect.Indirect(reflect.New(dest.Type().Elem()))
				err := unmarshalText(c, v, name+"[]", val)
				if err != nil {
					return err
				}
				m.SetMapIndex(reflect.ValueOf(key), v)
			}
			dest.Set(m)
		}
	default:
		return ctx.NewErrorf(nil, "unsupported field type: %T", v)
	}
	return nil
}
