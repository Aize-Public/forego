package storage

import (
	"sync"

	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
	"github.com/Aize-Public/forego/enc"
)

// A single table, where you can do random access or limited scans
type KeyValue interface {
	Get(c ctx.C, key string) (enc.Map, error)
	Upsert(c ctx.C, key string, val enc.Map) error
	Delete(c ctx.C, key string) error
	Range(c ctx.C, f func(c ctx.C, key string, val enc.Map) error, filters ...Filter) error
}

type memKeyValue struct {
	m    sync.Mutex
	data map[string]enc.Map
}

var _ KeyValue = NewMemKeyValue()

func NewMemKeyValue() *memKeyValue {
	return &memKeyValue{
		data: map[string]enc.Map{},
	}
}

func (this *memKeyValue) Get(c ctx.C, key string) (enc.Map, error) {
	this.m.Lock()
	defer this.m.Unlock()
	v := this.data[key]
	return v, nil
}

func (this *memKeyValue) Upsert(c ctx.C, key string, val enc.Map) error {
	this.m.Lock()
	defer this.m.Unlock()
	this.data[key] = val
	return nil
}

func (this *memKeyValue) Delete(c ctx.C, key string) error {
	this.m.Lock()
	defer this.m.Unlock()
	delete(this.data, key)
	return nil
}

func (this *memKeyValue) Range(c ctx.C, f func(c ctx.C, key string, val enc.Map) error, filters ...Filter) error {
	type pair struct {
		key string
		val enc.Map
	}
	var list []pair
	err := func() error {
		this.m.Lock()
		defer this.m.Unlock()
		for k, v := range this.data {
			ok, err := check(c, v, filters...)
			if err != nil {
				return err
			}
			log.Debugf(c, "prep %q is %v", k, ok)
			if ok {
				list = append(list, pair{k, v})
			}
		}
		return nil
	}()
	if err != nil {
		return err
	}
	for _, p := range list {
		err := f(c, p.key, p.val)
		if err == EOD {
			return nil
		} else if err != nil {
			return err
		}
	}
	return nil
}

func check(c ctx.C, val enc.Map, filters ...Filter) (bool, error) {
	for _, f := range filters {
		ok, err := f.Check(c, val)
		if err != nil || ok == false {
			return ok, err
		}
	}
	return true, nil
}

func (this Filter) Check(c ctx.C, val enc.Map) (bool, error) {
	v := val[this.Field]
	if v == nil {
		v = enc.Nil{}
	}
	log.Debugf(c, "Check(%q %q %v)", v, this.Cmp, this.Val)
	switch this.Cmp {
	case Equal:
		return v.GoString() == this.Val.GoString(), nil
	case NotEqual:
		return v.GoString() != this.Val.GoString(), nil
	default:
		return false, ctx.NewErrorf(c, "unsupported cmp operator: %q", this.Cmp)
	}
}
