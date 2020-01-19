package jsonrpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	// "github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"nhooyr.io/websocket"
)

var maxPacketSize = 1024 * 1024
var jsonProcessor = jsoniter.ConfigCompatibleWithStandardLibrary


func (s *Server) reader(ctx context.Context, conn *websocket.Conn, done chan struct{}, queue chan interface{}) {
	log := s.getlog()
	buffer := make([]byte, maxPacketSize)
	defer close(done)

	for {
		_, reader, err := conn.Reader(ctx)
		if err != nil {
			log.Error("wsClient.reader conn.Reader error: ", err)
			return
		}
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			log.Error("wsClient.reader Read(buffer) error: ", err)
			return
		}

		if n > 0 {
			tmp := make([]byte, n)
			copy(tmp, buffer)
			queue <- tmp
		}
	}
}

func (s *Server) writer(ctx context.Context, conn *websocket.Conn, done chan struct{}, queue chan interface{}) {
	log := s.getlog()
	defer close(done)
	writer := func(msg interface{}) error {
		w, err := conn.Writer(ctx, websocket.MessageText)
		if err != nil {
			return err
		}

		defer func() {
			if err := w.Close(); err != nil {
				log.Error("Server.writer writer error close: ", err)
			}
		}()
		e := jsonProcessor.NewEncoder(w)
		if err := e.Encode(msg);  err != nil {
			return fmt.Errorf("Server.writer failed to encode json: %w", err)
		}

		return nil
	}

	for {
		select {
		case msg := <-queue:
			s.submit(func() {
				if err := writer(msg); err != nil {
					log.Error("Server.writer write error: ", err)
				}
			})
		case <- ctx.Done():
			return
		}
	}
}

func (s *Server) Handle(ctx context.Context, conn *websocket.Conn) {
	log := s.getlog()
	var wg sync.WaitGroup
	var err error
	basectx, basecancel := context.WithCancel(ctx)

	readerStop := make(chan struct{})
	writerStop := make(chan struct{})

	readerqueue := make(chan interface{})
	writerqueue := make(chan interface{})
	responses := make(chan *Response)

	defer func() {
		wg.Wait()
		defer basecancel()
		close(responses)
		close(writerqueue)
		close(readerqueue)
	}()

	go s.reader(basectx, conn, readerStop, readerqueue)
	go s.writer(basectx, conn, writerStop, writerqueue)


	handleCtx := setupContext(ctx, responses)
	handleCtx, err = s.afterConnect(ctx)
	if err != nil {
		responses <- newResponseNotification("error", err.Error())
		return
	}

	for {
		select {
			case msg, _ := <- s.broadcast:
				writerqueue <- msg
			case rawreq := <- readerqueue:
				req := &Request{}
				if err := jsonProcessor.NewDecoder(bytes.NewReader(rawreq.([]byte))).Decode(req); err != nil {
					log.Error("Server.Handle result decoder error: ", err)
					return
				}
				wg.Add(1)
				s.submit(func() {
					defer wg.Done()
					rpcCallCtx, err := s.beforeRequest(handleCtx, req.Method, req.Params)
					if err != nil {
						responses <- newResponseError(req.ID, err.Error())
						return
					}

					method := s.methods[req.Method]
					if method == nil {
						responses <- handleNotFound(req)
						return
					}

					data := callMethod(rpcCallCtx, s.rcvr, method, req)

					responses <- data
				})
			case resp := <-responses:
				writerqueue <- resp
			case <- ctx.Done():
				return
			case <- readerStop:
				return
			case <- writerStop:
				return
		}
	}
}

func handleNotFound(req *Request) *Response {
	// log := s.getlog()
	rsp := newResponseError(req.ID, fmt.Errorf("method not found: %s", req.Method).Error())
	// log.Printf("rsp error: %s", rsp.Error)
	return rsp
}
