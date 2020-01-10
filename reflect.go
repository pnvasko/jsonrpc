package jsonrpc

import (
	"context"
	"reflect"
)

func convertParams(method *Method, req *Request) (interface{}, error) {
	if method.paramsType == nil {
		return nil, nil
	}
	params, err := req.Params.ParseInto(method.paramsType)
	if err != nil {
		return nil, err
	}
	return params, nil
}

func callMethod(ctx context.Context, t interface{}, method *Method, req *Request, params interface{}) *Response {
	in := []reflect.Value{
		reflect.ValueOf(t),
		reflect.ValueOf(ctx),
	}

	if method.paramsType != nil {
		in = append(in, reflect.ValueOf(params))
	}

	out := method.fn.Call(in)

	var err error
	var result interface{}
	switch len(out) {
	case 0:
	case 1:
		err = getError(out[0])
	case 2:
		result = getResult(out[0])
		err = getError(out[1])
	default:
		panic("invalid # of arguments")
	}

	if err != nil {
		if userErr, ok := err.(interface{ UserError() string }); ok {
			return newResponseError(req.ID, userErr.UserError())
		}
		return newResponseError(req.ID, err.Error())
	}
	return newResponse(req.ID, result)
}

func getResult(out reflect.Value) interface{} {
	if out.Kind() != reflect.Ptr {
		return out.Interface()
	}
	if out.IsNil() {
		return nil
	}
	return out.Interface()
}

func getError(out reflect.Value) error {
	err, _ := getResult(out).(error)
	return err
}
