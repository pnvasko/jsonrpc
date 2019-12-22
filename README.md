A [JSON-RPC 2.0](https://www.jsonrpc.org/specification) implementation for Go.

Supports any transport that can pass strings back and forth like WebSockets.

Example:

```go
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

type MyFunctionParams struct{
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
```

Usage with [websocat](https://github.com/vi/websocat):

```sh-session
$ websocat --jsonrpc ws://localhost:8000/rpc
{"id":1,"method":"ExampleFunc","params":{"should_error": false}}
{"jsonrpc":"2.0","id":1,"result":"result can be anything json-marshalable"}
{"id":2,"method":"ExampleFunc","params":{"should_error": true}}
{"jsonrpc":"2.0","id":2,"error":"this error returned to client"}

# or use the --jsonrpc flag
$ websocat --jsonrpc ws://localhost:8000/rpc
ExampleFunc {}
{"jsonrpc":"2.0","id":1,"result":"result can be anything json-marshalable"}
ExampleFunc {"should_error": true}
{"jsonrpc":"2.0","id":2,"error":"this error returned to client"}
```
