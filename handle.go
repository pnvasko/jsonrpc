package jsonrpc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"runtime/debug"
)

type Socket interface {
	ReadJSON(interface{}) error
	WriteJSON(interface{}) error
	Close() error
}

func readRequests(sock Socket) <-chan *Request {
	requestChan := make(chan *Request)
	go func() {
		defer close(requestChan)
		for {
			var req Request
			if err := sock.ReadJSON(&req); err != nil {
				log.Printf("req error: %+v", err)
				return
			}
			log.Printf("req: %d %s", req.ID, req.Method)
			requestChan <- &req
		}
	}()
	return requestChan
}

func callMethod(ctx context.Context, t interface{}, method *Method, req *Request) *Response {
	in := []reflect.Value{
		reflect.ValueOf(t),
		reflect.ValueOf(ctx),
	}

	if method.paramsType != nil {
		params, err := req.Params.ParseInto(method.paramsType)
		if err != nil {
			return newResponseError(req.ID, err.Error())
		}
		log.Printf("req: %d %s %+v", req.ID, req.Method, params)

		in = append(in, reflect.ValueOf(params))
	}

	out := method.fn.Call(in)

	var err error
	var result interface{}
	switch len(out) {
	case 0:
	case 1:
		err = getError(out[0])
	case 2:
		result = getResult(out[0])
		err = getError(out[1])
	default:
		panic("invalid # of arguments")
	}

	if err != nil {
		if userErr, ok := err.(interface{ UserError() string }); ok {
			return newResponseError(req.ID, userErr.UserError())
		}
		return newResponseError(req.ID, err.Error())
	}
	return newResponse(req.ID, result)
}

func handleNotFound(req *Request) *Response {
	rsp := newResponseError(req.ID, fmt.Errorf("method not found: %s", req.Method).Error())
	log.Printf("rsp error: %s", rsp.Error)
	return rsp
}

func handlePanic(req *Request, responses chan<- *Response) {
	errish := recover()
	if errish == nil {
		return
	}
	rsp := newResponseError(req.ID, errors.New("internal server error").Error())
	log.Printf("%+v", errish)

	// TODO: hide error in production
	rsp.Error = fmt.Sprintf("%+v", errish)

	responses <- rsp
}

func getResult(out reflect.Value) interface{} {
	if out.Kind() != reflect.Ptr {
		return out.Interface()
	}
	if out.IsNil() {
		return nil
	}
	return out.Interface()
}

func getError(out reflect.Value) error {
	err, _ := getResult(out).(error)
	return err
}

func writeResponses(sock Socket, responses <-chan *Response) {
	for rsp := range responses {
		if rsp.Error != "" {
			log.Printf("rsp error: %d %s", rsp.ID, rsp.Error)
		} else {
			log.Printf("rsp: %d", rsp.ID)
		}
		debug.PrintStack()
		if err := sock.WriteJSON(rsp); err != nil {
			log.Println(err)
		}
	}
}
