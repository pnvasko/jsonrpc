package jsonrpc

import (
	"context"
)

func onClose(sock Socket, responses chan<- *Response) {
	close(responses)
	sock.Close()
}

func Close(ctx context.Context) {
	ctxGetCloseFunc(ctx)()
}
