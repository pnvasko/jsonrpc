package jsonrpc

type Response struct {
	ID      ID          `json:"id,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   string      `json:"error,omitempty"`
	JSONRPC string      `json:"jsonrpc"`
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
