package jsonrpc

import (
	"context"
	"sync"
)

func onClose(sock Socket, responses chan<- *Response, wg *sync.WaitGroup) {
	wg.Wait()
	close(responses)
	sock.Close()
}

func Close(ctx context.Context) {
	ctxGetCloseFunc(ctx)()
}
