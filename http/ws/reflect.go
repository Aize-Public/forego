package ws

import (
	"reflect"
	"strings"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
)

// used by wsrpc

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
		//log.Debugf(nil, "set %v", c.value)
		v.Elem().Field(c.index).Set(c.value)
	}

	if this.constructor.methodName != "" {
		err := this.constructor.call(c, v, req)
		if err != nil {
			return err
		}
	}

	// setup channel routing...
	c.ch.byPath = map[string]func(c C, req enc.Node) error{}
	for _, method := range this.methods {
		method := method
		c.ch.byPath[method.name] = func(c C, req enc.Node) error {
			err := method.call(c, v, req)
			if err != nil {
				c.ch.Conn.Send(c, Frame{
					Channel: c.ch.ID,
					Path:    method.name,
					Type:    "return",
					Data:    enc.MustMarshal(c, err),
				})
				return ctx.NewErrorf(c, "ws[%s|%s]: %v", c.ch.ID, method.name, err)
			}
			c.ch.Conn.Send(c, Frame{
				Channel: c.ch.ID,
				Path:    method.name,
				Type:    "return",
			})
			return nil
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
			return ctx.NewErrorf(c, "reflect[%v].argument: %w", this.name, err)
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

func (this *builder) inspect(c ctx.C, obj any) error {
	log.Debugf(c, "ws.inspect(%T %v)", obj, obj)
	origVal := reflect.ValueOf(obj)
	switch origVal.Kind() {
	case reflect.Pointer:
	case reflect.Struct:
		origVal = origVal.Addr() // make sure it's a pointer
	default:
		return ctx.NewErrorf(c, "expected struct or *struct, got %T", obj)
	}
	this.structType = origVal.Type().Elem()
	ptrType := origVal.Type()

	this.name = toLowerFirst(this.structType.Name())
	log.Infof(c, "WS object %q: %v", this.name, this.structType)

	// shallow copy fields value to the new obj
	for i := 0; i < origVal.Elem().NumField(); i++ {
		i := i
		fv := origVal.Elem().Field(i)
		if !fv.IsZero() {
			this.fields = append(this.fields, fieldInit{
				index: i,
				value: fv,
			})
		}
	}

	// scan methods
	for i := 0; i < ptrType.NumMethod(); i++ {
		m := ptrType.Method(i)
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
			this.constructor = method
		} else {
			this.methods = append(this.methods, method)
		}
	}
	return nil
}

func toLowerFirst(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToLower(s[0:1]) + s[1:]
}
