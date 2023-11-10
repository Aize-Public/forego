package ws

import (
	"reflect"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
)

type builder struct {
	name        string
	structType  reflect.Type
	fields      []fieldInit
	constructor method
	methods     []method
}

// init a new object, and bind it's methods to the channel
func (this builder) build(c C, req enc.Node) any {
	v := reflect.New(this.structType)
	for _, c := range this.fields {
		log.Debugf(nil, "set %v", c.value)
		v.Elem().Field(c.index).Set(c.value)
	}

	if this.constructor.methodName != "" {
		this.constructor.call(c, v, req)
	}

	// setup channel routing...
	c.ch.byPath = map[string]func(c C, req enc.Node) error{}
	for _, method := range this.methods {
		c.ch.byPath[method.name] = func(c C, req enc.Node) error {
			return method.call(c, v, req)
		}
	}

	return v.Interface()
}

type method struct {
	name       string
	methodName string
	argument   reflect.Type // may be nil
}

func (this method) call(c C, obj reflect.Value, request enc.Node) error {
	args := []reflect.Value{
		reflect.ValueOf(c),
	}
	if this.argument != nil {
		inv := reflect.New(this.argument)
		err := enc.Unmarshal(c, request, inv.Interface())
		if err != nil {
			return err
		}
		args = append(args, inv.Elem())
	}
	m := obj.MethodByName(this.methodName)
	ret := m.Call(args)
	switch len(ret) {
	case 0:
		return nil
	default:
		// assume the last one is error
		ev := ret[len(ret)-1]
		if ev.IsNil() {
			return nil
		}
		return ev.Interface().(error)
	}
}

type fieldInit struct {
	index int
	value reflect.Value
}

func inspect(c ctx.C, obj any) (builder, error) {
	rv := reflect.ValueOf(obj)
	rt := rv.Type()
	st := rt
	if st.Kind() == reflect.Pointer {
		st = rt.Elem()
	} else {
		return builder{}, ctx.NewErrorf(c, "must be a pointer to be reasonable state: %T", obj)
	}

	builder := builder{
		name:       toLowerFirst(st.Name()),
		structType: st,
	}
	log.Infof(c, "WS object %q: %v", builder.name, builder.structType)

	// shallow copy fields value to the new obj
	for i := 0; i < rv.Elem().NumField(); i++ {
		fv := rv.Elem().Field(i)
		if !fv.IsZero() {
			builder.fields = append(builder.fields, fieldInit{
				index: i,
				value: fv,
			})
		}
	}

	// scan methods
	for i := 0; i < rt.NumMethod(); i++ {
		m := rt.Method(i)
		method := method{
			name:       toLowerFirst(m.Name),
			methodName: m.Name,
		}
		switch m.Type.NumIn() {
		case 3:
			method.argument = m.Type.In(2)
			if m.Type.In(1) != reflect.TypeOf(C{}) {
				log.Debugf(c, "WS ignoring %q because first arg is %v", method.name, method.argument)
				continue
			} else {
				log.Infof(c, "WS handler %q with argument %v", method.name, method.argument)
			}
		case 2:
			if m.Type.In(1) != reflect.TypeOf(C{}) {
				log.Debugf(c, "WS ignoring %q because first arg is %v", method.name, method.argument)
				continue
			} else {
				log.Infof(c, "WS handler %q with no args", method.name)
			}
		default:
			log.Debugf(c, "WS ignoring %q because %d args", method.name, method.argument)
			continue
		}
		if m.Name == "Init" {
			builder.constructor = method
		} else {
			builder.methods = append(builder.methods, method)
		}
	}
	return builder, nil
}
