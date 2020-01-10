package jsonrpc

import (
	"context"
)

func Notify(ctx context.Context, method string, params interface{}) {
	ctxGetNotifyFunc(ctx)(method, params)
}
