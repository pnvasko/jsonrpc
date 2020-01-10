package jsonrpc

import (
	"io"
	"reflect"
)

type Server struct {
	methods       Methods
	rcvr          interface{}
	afterConnect  afterConnectFN
	beforeRequest beforeRequestFN
}

type Socket interface {
	io.Closer
	ReadJSON(interface{}) error
	WriteJSON(interface{}) error
}

func New(sampleMethodReceiver interface{}) *Server {
	methods := Methods{}
	ty := reflect.TypeOf(sampleMethodReceiver)
	for i := 0; i < ty.NumMethod(); i++ {
		m := ty.Method(i)
		fn := m.Func
		fnType := fn.Type()
		var paramsType reflect.Type
		if fnType.NumIn() == 3 {
			paramsType = fnType.In(2)
		}
		methods[m.Name] = &Method{fn, paramsType}
	}

	return &Server{
		methods:       methods,
		rcvr:          sampleMethodReceiver,
		afterConnect:  getAfterConnect(sampleMethodReceiver),
		beforeRequest: getBeforeRequest(sampleMethodReceiver),
	}
}
