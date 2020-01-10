package jsonrpc

import (
	"context"
)

type (
	ctxNotifyFuncKey struct{}
	ctxCloseFuncKey  struct{}
)

func ctxGetNotifyFunc(ctx context.Context) func(method string, params interface{}) {
	return ctx.Value(ctxNotifyFuncKey{}).(func(method string, params interface{}))
}

func ctxWithNotifyFunc(ctx context.Context, fn func(method string, params interface{})) context.Context {
	return context.WithValue(ctx, ctxNotifyFuncKey{}, fn)
}

func ctxGetCloseFunc(ctx context.Context) func() {
	return ctx.Value(ctxCloseFuncKey{}).(func())
}

func ctxWithCloseFunc(ctx context.Context, fn func()) context.Context {
	return context.WithValue(ctx, ctxCloseFuncKey{}, fn)
}

func setupContext(ctx context.Context, responses chan<- *Response) context.Context {
	ctx, cancel := context.WithCancel(ctx)
	ctx = ctxWithCloseFunc(ctx, cancel)
	ctx = ctxWithNotifyFunc(ctx, func(method string, params interface{}) {
		responses <- newResponseNotification(method, params)
	})
	return ctx
}
