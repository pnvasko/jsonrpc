package jsonrpc

import (
	"context"
)

type BeforeRequest interface {
	BeforeRequest(ctx context.Context, method string, params interface{}) (context.Context, error)
}

func getBeforeMiddleware(rcvr interface{}) func(
	ctx context.Context, method string, params interface{}) (context.Context, error) {
	r, ok := rcvr.(BeforeRequest)
	if !ok {
		return beforeRequestNoop
	}
	return r.BeforeRequest
}

func beforeRequestNoop(ctx context.Context, method string, params interface{}) (
	context.Context, error) {
	return ctx, nil
}
