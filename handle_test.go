package jsonrpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"github.com/jdxcode/jsonrpc"
)

var rpc = jsonrpc.New(&TestRPC{})

func TestHandleString(t *testing.T) {
	assert := assert.New(t)
	sock := newFakeSocket()
	go func() {
		defer close(sock.requests)
		params := jsonrpc.ParamsRaw("\"test-abc\"")
		sock.requests <- &jsonrpc.Request{ID: 101, Method: "Foo", Params: &params}
		rsp := <-sock.responses
		assert.Equal(jsonrpc.ID(101), rsp.ID)
		assert.Equal(123, rsp.Result)
		assert.Empty(rsp.Error)
	}()
	rpc.NewSession(&TestRPC{}, sock).HandleRequests(ctx)
}

func TestHandleStruct(t *testing.T) {
	assert := assert.New(t)
	sock := newFakeSocket()
	go func() {
		defer close(sock.requests)
		params := jsonrpc.ParamsRaw("{\"foo\": \"test-abc\"}")
		sock.requests <- &jsonrpc.Request{ID: 102, Method: "FooStruct", Params: &params}
		rsp := <-sock.responses
		assert.Equal(jsonrpc.ID(102), rsp.ID)
		result := rsp.Result.(*FooStructResult)
		assert.Equal("test-abc", result.Bar)
		assert.Empty(rsp.Error)
	}()
	rpc.NewSession(&TestRPC{}, sock).HandleRequests(ctx)
}

func TestHandleErr(t *testing.T) {
	assert := assert.New(t)
	sock := newFakeSocket()
	go func() {
		defer close(sock.requests)
		params := jsonrpc.ParamsRaw("\"test-abc\"")
		sock.requests <- &jsonrpc.Request{ID: 102, Method: "FooErr", Params: &params}
		rsp := <-sock.responses
		assert.Equal(jsonrpc.ID(102), rsp.ID)
		assert.Nil(rsp.Result)
		assert.Equal("uh oh", rsp.Error)
	}()
	rpc.NewSession(&TestRPC{}, sock).HandleRequests(ctx)
}

func TestHandlePanic(t *testing.T) {
	assert := assert.New(t)
	sock := newFakeSocket()
	go func() {
		defer close(sock.requests)
		sock.requests <- &jsonrpc.Request{ID: 102, Method: "FooPanic"}
		rsp := <-sock.responses
		assert.Equal(jsonrpc.ID(102), rsp.ID)
		assert.Nil(rsp.Result)
		assert.Equal("uh oh", rsp.Error)
	}()
	rpc.NewSession(&TestRPC{}, sock).HandleRequests(ctx)
}

func TestMethodNotFound(t *testing.T) {
	assert := assert.New(t)
	sock := newFakeSocket()
	go func() {
		sock.requests <- &jsonrpc.Request{ID: 101, Method: "invalid_method", Params: nil}
		rsp := <-sock.responses
		assert.Equal(jsonrpc.ID(101), rsp.ID)
		assert.Nil(rsp.Result)
		assert.Equal("method not found: invalid_method", rsp.Error)
		close(sock.requests)
	}()
	rpc.NewSession(&TestRPC{}, sock).HandleRequests(ctx)
}

type FakeSocket struct {
	requests  chan *jsonrpc.Request
	responses chan *jsonrpc.Response
}

func newFakeSocket() *FakeSocket {
	return &FakeSocket{make(chan *jsonrpc.Request), make(chan *jsonrpc.Response)}
}

func (f *FakeSocket) ReadJSON(raw interface{}) error {
	cur := <-f.requests
	if cur == nil {
		return &websocket.CloseError{
			Code: websocket.CloseNormalClosure,
			Text: "closed normally",
		}
	}
	request := raw.(*jsonrpc.Request)
	request.ID = cur.ID
	request.Method = cur.Method
	request.Params = cur.Params
	return nil
}

func (f *FakeSocket) WriteJSON(response interface{}) error {
	f.responses <- response.(*jsonrpc.Response)
	return nil
}

func (f *FakeSocket) Close() error {
	return nil
}

type TestRPC struct{}

func (r *TestRPC) Foo(ctx context.Context, params string) (int, error) {
	return 123, nil
}

type FooStructParams struct {
	Foo string `json:"foo"`
}

type FooStructResult struct {
	Bar string
}

func (r *TestRPC) FooStruct(
	ctx context.Context, params *FooStructParams) (*FooStructResult, error) {
	return &FooStructResult{params.Foo}, nil
}

func (r *TestRPC) FooErr(ctx context.Context, params string) (interface{}, error) {
	return nil, errors.New("uh oh")
}

func (r *TestRPC) FooPanic(ctx context.Context) (interface{}, error) {
	panic("uh oh")
}

var ctx = context.Background()
