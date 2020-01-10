package jsonrpc

import (
	"context"
)

type AfterConnect interface {
	AfterConnect(ctx context.Context) (context.Context, error)
}

type BeforeRequest interface {
	BeforeRequest(ctx context.Context, method string, params interface{}) (context.Context, error)
}

func getAfterConnect(rcvr interface{}) func(ctx context.Context) (context.Context, error) {
	r, ok := rcvr.(AfterConnect)
	if !ok {
		return afterConnectNoop
	}
	return r.AfterConnect
}

func afterConnectNoop(ctx context.Context) (context.Context, error) {
	return ctx, nil
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
