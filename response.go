package jsonrpc

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

func newResponseError(id ID, err error) *Response {
	return &Response{
		ID:      id,
		Error:   err.Error(),
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
