package jsonrpc

import (
	"context"
	"io"
	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/pnvasko/ants"
)

type Server struct {
	ctx 		  context.Context
	pool      	  *ants.Pool
	broadcast	  chan *Response

	methods       Methods
	rcvr          interface{}
	afterConnect  afterConnectFN
	beforeRequest beforeRequestFN

	controlRpc	  func(interface{}) interface{}
}

type Socket interface {
	io.Closer
	ReadJSON(interface{}) error
	WriteJSON(interface{}) error
}

func New(ctx context.Context, cf func(interface{}) interface{}, sampleMethodReceiver interface{}) *Server {
	var err error
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

	srv := &Server{
		ctx: 		   ctx,
		broadcast:	   make(chan *Response),
		methods:       methods,
		rcvr:          sampleMethodReceiver,
		afterConnect:  getAfterConnect(sampleMethodReceiver),
		beforeRequest: getBeforeRequest(sampleMethodReceiver),
		controlRpc:	   cf,
	}
	log := srv.getlog()

	if srv.pool, err = ants.NewPool(ctx, 100, ants.WithPreAlloc(true)); err != nil {
		log.Error("app_pool.NewPool error: ", err)
	}

	return srv
}

func (s *Server) getlog() *logrus.Logger {
	return s.ctx.Value("log").(*logrus.Logger)
}

func (s *Server) submit(f func()) {
	log := s.getlog()
	if err := s.pool.Submit(f); err != nil {
		log.Error("Server.submit error: ", err)
	}
}
/*
func (s *Server) Call(method string, params interface{}) (interface{}, error) {
	basectx, basecancel := context.WithCancel(s.ctx)
	responses := make(chan *Response)
	defer func() {
		defer basecancel()
		close(responses)
	}()
	handleCtx := setupContext(s.ctx, responses)
	handleCtx, err := s.afterConnect(handleCtx)
	if err != nil {
		responses <- newResponseNotification("error", err.Error())
		return nil, err
	}

	rpcCallCtx, err := s.beforeRequest(handleCtx, method, params)
	if err != nil {
		responses <- newResponseError(req.ID, err.Error())
		return
	}

	resp := <-responses

	return resp, nil
}
*/
func (s *Server) Broadcast() chan <- *Response {
	return s.broadcast
}

func (s *Server) Close()  {
	// close(s.broadcast)
}