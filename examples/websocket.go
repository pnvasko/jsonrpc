package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/jdxcode/jsonrpc"
)

const PORT = ":8000"

// starts a json-rpc 2.0 websocket server
func main() {
	rpc := jsonrpc.New(&RPC{})
	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		upgrader := &websocket.Upgrader{}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		ctx := r.Context()
		rpc.Handle(ctx, conn)
	})

	if err := http.ListenAndServe(PORT, nil); err != nil {
		log.Fatal(err)
	}
}

type RPC struct{}

type MyFunctionParams struct {
	ShouldError bool `json:"should_error"`
}

// the client will call this by sending a message like:
// {"jsonrpc": "2.0", id: 101, method: "ExampleFunc", params: {"should_error": false}}
// the server would respond:
// {"jsonrpc": "2.0", id: 101, result: "result can be anything json-marshalable"}
func (r *RPC) ExampleFunc(ctx context.Context, params *MyFunctionParams) (string, error) {
	if params.ShouldError {
		return "", errors.New("this error returned to client")
	}
	return "result can be anything json-marshalable", nil
}

func (r *RPC) NoParamsOrResult(ctx context.Context) error {
	return nil
}
