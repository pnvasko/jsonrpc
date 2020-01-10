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
			defer wg.Done()
			defer handlePanic(req, responses)

			var err error
			ctx, err = s.beforeRequest(ctx, req.Method, req.Params)
			if err != nil {
				responses <- newResponseError(req.ID, err.Error())
				wg.Done()
				return
			}

			method := s.methods[req.Method]
			if method == nil {
				responses <- handleNotFound(req)
				return
			}
			responses <- callMethod(ctx, s.rcvr, method, req)
		}(req)
	}
}

func handleNotFound(req *Request) *Response {
	rsp := newResponseError(req.ID, fmt.Errorf("method not found: %s", req.Method).Error())
	log.Printf("rsp error: %s", rsp.Error)
	return rsp
}
