package jsonrpc

import (
	"reflect"
)

type Server struct {
	methods Methods
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
		methods: methods,
	}
}

func (s *Server) NewSession(rcvr interface{}, sock Socket) *Session {
	responses := make(chan *Response)
	go writeResponses(sock, responses)
	return &Session{
		rcvr:      rcvr,
		sock:      sock,
		methods:   s.methods,
		responses: responses,
	}
}
