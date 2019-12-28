package jsonrpc

import (
	"reflect"
)

type Methods map[string]*Method

type Method struct {
	fn         reflect.Value
	paramsType reflect.Type
}
