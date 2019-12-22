package jsonrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Request struct {
	JSONRPC string     `json:"jsonrpc"`
	ID      ID         `json:"id"`
	Method  string     `json:"method"`
	Params  *ParamsRaw `json:"params"`
}

type (
	ID        int
	ParamsRaw []byte
)

func (p *ParamsRaw) UnmarshalJSON(b []byte) error {
	*p = b
	return nil
}

func (p *ParamsRaw) ParseInto(paramsType reflect.Type) (interface{}, error) {
	var params interface{}
	if paramsType.Kind() == reflect.Ptr {
		params = reflect.New(paramsType.Elem()).Interface()
	} else {
		params = reflect.New(paramsType).Elem().Interface()
	}
	if err := json.Unmarshal(*p, &params); err != nil {
		return nil, fmt.Errorf("rpc [params unmarshal]: %w", err)
	}
	return params, nil
}
