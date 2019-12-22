package jsonrpc

import (
	"reflect"
)

type JSONRPC struct {
	t       interface{}
	methods map[string]*method
}

type method struct {
	fn         reflect.Value
	paramsType reflect.Type
}

func New(t interface{}) *JSONRPC {
	methods := map[string]*method{}

	ty := reflect.TypeOf(t)
	for i := 0; i < ty.NumMethod(); i++ {
		m := ty.Method(i)
		fn := m.Func
		fnType := fn.Type()
		var paramsType reflect.Type
		if fnType.NumIn() == 3 {
			paramsType = fnType.In(2)
		}
		methods[m.Name] = &method{fn, paramsType}
	}

	return &JSONRPC{
		t:       t,
		methods: methods,
	}
}
