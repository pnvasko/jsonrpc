package jsonrpc

func (s *Session) Notify(method string, params interface{}) {
	s.responses <- newResponseNotification(method, params)
}
