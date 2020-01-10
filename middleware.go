package jsonrpc

import (
	"context"
)

type afterConnectFN = func(ctx context.Context) (context.Context, error)

type AfterConnect interface {
	AfterConnect(ctx context.Context) (context.Context, error)
}

type (
	beforeRequestFN = func(ctx context.Context, method string, params interface{}) (context.Context, error)
	BeforeRequest   interface {
		BeforeRequest(ctx context.Context, method string, params interface{}) (context.Context, error)
	}
)

func getAfterConnect(rcvr interface{}) afterConnectFN {
	r, ok := rcvr.(AfterConnect)
	if !ok {
		return afterConnectNoop
	}
	return r.AfterConnect
}

func afterConnectNoop(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func getBeforeRequest(rcvr interface{}) beforeRequestFN {
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
