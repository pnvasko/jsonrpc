package jsonrpc

import (
	"context"
	"fmt"
	"log"
	"sync"
)

func (s *Server) Handle(ctx context.Context, sock Socket) {
	var err error
	var wg sync.WaitGroup

	responses := make(chan *Response)
	go writeResponses(sock, responses)
	defer onClose(sock, responses, &wg)

	ctx = setupContext(ctx, responses)

	ctx, err = s.afterConnect(ctx)
	if err != nil {
		responses <- newResponseNotification("error", err.Error())
		return
	}

	for req := range readRequests(ctx, sock) {
		wg.Add(1)
		go func(req *Request) {
			var err error
			defer wg.Done()
			defer handlePanic(req, responses)

			method := s.methods[req.Method]
			if method == nil {
				responses <- handleNotFound(req)
				return
			}
			params, err := convertParams(method, req)
			if err != nil {
				responses <- newResponseError(req.ID, err.Error())
			}
			log.Printf("req: %d %s %+v", req.ID, req.Method, params)

			ctx, err = s.beforeRequest(ctx, req.Method, params)
			if err != nil {
				responses <- newResponseError(req.ID, err.Error())
				return
			}

			responses <- callMethod(ctx, s.rcvr, method, req, params)
		}(req)
	}
}

func handleNotFound(req *Request) *Response {
	rsp := newResponseError(req.ID, fmt.Errorf("method not found: %s", req.Method).Error())
	log.Printf("rsp error: %s", rsp.Error)
	return rsp
}
