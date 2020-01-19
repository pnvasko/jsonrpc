package jsonrpc

import (
	"log"
)

type Response struct {
	ID     ID          `json:"id,omitempty"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`

	// for notifications to client
	Method string      `json:"method,omitempty"`
	Params interface{} `json:"params,omitempty"`

	JSONRPC string `json:"jsonrpc"`
}

func newResponse(id ID, result interface{}) *Response {
	return &Response{
		ID:      id,
		Result:  result,
		JSONRPC: "2.0",
	}
}
// TODO make by standart
func newResponseError(id ID, err string) *Response {
	return &Response{
		ID:      id,
		Error:   err,
		JSONRPC: "2.0",
	}
}

func newResponseNotification(method string, params interface{}) *Response {
	return &Response{
		Method:  method,
		Params:  params,
		JSONRPC: "2.0",
	}
}

func writeResponses(sock Socket, responses <-chan *Response) {
	for rsp := range responses {
		if rsp.Error != "" {
			log.Printf("rsp error: %d %s", rsp.ID, rsp.Error)
		} else {
			log.Printf("rsp: %d", rsp.ID)
		}
		if err := sock.WriteJSON(rsp); err != nil {
			log.Println(err)
		}
	}
}
