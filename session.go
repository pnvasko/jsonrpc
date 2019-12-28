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
	for req := range readRequests(s.sock) {
		go func(req *Request) {
			defer handlePanic(req, s.responses)

			method := s.methods[req.Method]
			if method == nil {
				s.responses <- handleNotFound(req)
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
