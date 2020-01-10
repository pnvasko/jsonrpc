package jsonrpc

import (
	"errors"
	"fmt"
	"log"
	"runtime/debug"
)

func handlePanic(req *Request, responses chan<- *Response) {
	errish := recover()
	if errish == nil {
		return
	}
	debug.PrintStack()
	rsp := newResponseError(req.ID, errors.New("internal server error").Error())
	log.Printf("%+v", errish)

	// TODO: hide error in production
	rsp.Error = fmt.Sprintf("%+v", errish)

	responses <- rsp
}
