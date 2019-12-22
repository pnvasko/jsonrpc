package jsonrpc

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      ID          `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func newResponse(id ID, result interface{}) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

func newResponseError(id ID, err error) *Response {
	return &Response{
		JSONRPC: "2.0",
		ID:      id,
		Error:   err.Error(),
	}
}
