package jsonrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type Request struct {
	ID      ID         `json:"id"`
	Method  string     `json:"method"`
	Params  *ParamsRaw `json:"params"`
	JSONRPC string     `json:"jsonrpc"`
}

type (
	ID        int
	ParamsRaw []byte
)

func (p *ParamsRaw) UnmarshalJSON(b []byte) error {
	*p = b
	return nil
}

func (p *ParamsRaw) ParseInto(paramsType reflect.Type) (interface{}, error) {
	var params interface{}
	if paramsType.Kind() == reflect.Ptr {
		params = reflect.New(paramsType.Elem()).Interface()
	} else {
		params = reflect.New(paramsType).Elem().Interface()
	}
	if err := json.Unmarshal(*p, &params); err != nil {
		return nil, fmt.Errorf("rpc [params unmarshal]: %w", err)
	}
	return params, nil
}

func readRequests(ctx context.Context, sock Socket) <-chan *Request {
	requests := make(chan *Request)
	go func() {
		defer close(requests)
		for {
			select {
			case r := <-readNextRequest(sock):
				if r.err != nil {
					log.Printf("req error: %+v", r.err)
					return
				}
				log.Printf("req: %d %s", r.req.ID, r.req.Method)
				requests <- r.req
			case <-ctx.Done():
				return
			}
		}
	}()
	return requests
}

type nextRequestResult struct {
	req *Request
	err error
}

func readNextRequest(sock Socket) <-chan nextRequestResult {
	ch := make(chan nextRequestResult)
	go func() {
		var req Request
		if err := sock.ReadJSON(&req); err != nil {
			ch <- nextRequestResult{nil, err}
			return
		}
		ch <- nextRequestResult{&req, nil}
	}()
	return ch
}
