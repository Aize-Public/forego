package ws

import (
	"reflect"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
)

func inspect[T any](c ctx.C, obj T) (string, func(ch *Channel, data enc.Node) T, error) {
	rv := reflect.ValueOf(obj)
	rt := rv.Type()
	st := rt
	if st.Kind() == reflect.Pointer {
		st = rt.Elem()
	} else {
		return "", nil, ctx.NewErrorf(c, "must be a pointer to be reasonable state: %T", obj)
	}

	name := toLowerFirst(st.Name())
	log.Infof(c, "WS object %q: %v", name, rt)

	var init []func(v reflect.Value)
	for i := 0; i < rv.Elem().NumField(); i++ {
		i := i
		fv := rv.Elem().Field(i)
		if !fv.IsZero() {
			init = append(init, func(v reflect.Value) {
				log.Debugf(nil, "set %v", fv)
				v.Elem().Field(i).Set(fv)
			})
		}
	}

	byPath := map[string]func(c C, v reflect.Value, n enc.Node) error{}
	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		mname := toLowerFirst(m.Name)
		switch m.Type.NumIn() {
		case 3:
			in := m.Type.In(2)
			if m.Type.In(1) != reflect.TypeOf(C{}) {
				log.Debugf(c, "WS ignoring %q because first arg is %v", mname, m.Type.In(1))
			} else {
				log.Infof(c, "WS handler %q with argument %v", mname, in)
				byPath[mname] = func(c C, v reflect.Value, n enc.Node) error {
					//log.Debugf(c, "call %v.%s(%v %v)", rt, mname, in, n)
					inv := reflect.New(in)
					err := enc.Unmarshal(c, n, inv.Interface())
					if err != nil {
						return err
					}
					vals := v.MethodByName(m.Name).Call([]reflect.Value{
						reflect.ValueOf(c),
						inv.Elem(),
					})
					switch len(vals) {
					case 0:
						return nil
					default:
						// assume the last one is error
						ev := vals[len(vals)-1]
						if ev.IsNil() {
							return nil
						}
						return ev.Interface().(error)
					}
				}
			}
		case 2:
			if m.Type.In(1) != reflect.TypeOf(C{}) {
				log.Debugf(c, "WS ignoring %q because first arg is %v", mname, m.Type.In(1))
			} else {
				log.Infof(c, "WS handler %q with no args", mname)
				byPath[mname] = func(c C, v reflect.Value, n enc.Node) error {
					//log.Debugf(c, "ch(%q) %v <= %v", mname, v, n)
					vals := v.MethodByName(m.Name).Call([]reflect.Value{
						reflect.ValueOf(c),
					})
					switch len(vals) {
					case 0:
						return nil
					default:
						// assume the last one is error
						ev := vals[len(vals)-1]
						if ev.IsNil() {
							return nil
						}
						return ev.Interface().(error)
					}
				}
			}
		default:
			log.Debugf(c, "WS ignoring %q because %d args", mname, m.Type.NumIn())
		}
	}
	return name, func(ch *Channel, data enc.Node) T {
		v := reflect.New(st)
		for _, f := range init {
			f(v)
		}
		ch.byPath = map[string]func(c C, n enc.Node) error{}
		for path, h := range byPath {
			path := path
			h := h
			if path == "init" {
				log.Debugf(c, "WS init(%v)", data)
				h(C{
					C:  c,
					ch: ch,
				}, v, data)
			} else {
				ch.byPath[path] = func(c C, n enc.Node) error {
					log.Debugf(c, "WS reflecting %+v .%s", v, path)
					return h(c, v, n)
				}
			}
		}
		return v.Interface().(T)
	}, nil
}
