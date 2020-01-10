package jsonrpc

import (
	"context"
)

type Session struct {
	methods Methods
	rcvr    interface{}
	sock    Socket

	responses chan<- *Response
}

func (s *Session) HandleRequests(ctx context.Context) {
	defer s.Close()
	beforeMiddleware := getBeforeMiddleware(s.rcvr)
	for req := range readRequests(s.sock) {
		go func(req *Request) {
			defer handlePanic(req, s.responses)

			method := s.methods[req.Method]
			if method == nil {
				s.responses <- handleNotFound(req)
				return
			}
			var err error
			ctx, err = beforeMiddleware(ctx, req.Method, req.Params)
			if err != nil {
				s.responses <- newResponseError(req.ID, err.Error())
				return
			}
			s.responses <- callMethod(ctx, s.rcvr, method, req)
		}(req)
	}
}

func (s *Session) Close() {
	s.sock.Close()
	close(s.responses)
}
