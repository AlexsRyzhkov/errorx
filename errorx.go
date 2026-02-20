package errorx

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"go.uber.org/zap/zapcore"
)

type errorStruct struct {
	Namespace string
	Code      string
	Message   string
	Fields    map[string]any
	Caller    string
	Type      string
	Err       error
}

func (e *errorStruct) Error() string {
	return e.Message
}

func (e *errorStruct) Unwrap() error {
	return e.Err
}

func (e *errorStruct) MarshalJSON() ([]byte, error) {

	return json.Marshal(e.getChain())
}

func (e *errorStruct) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for _, err := range e.getChain() {
		_ = arr.AppendObject(err)
	}
	return nil
}

type node struct {
	Namespace string         `json:"namespace,omitempty"`
	Code      string         `json:"code,omitempty"`
	Message   string         `json:"message,omitempty"`
	Fields    map[string]any `json:"fields,omitempty"`
	Caller    string         `json:"caller,omitempty"`
	Type      string         `json:"type,omitempty"`
}

func (n node) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if n.Namespace != "" {
		enc.AddString("namespace", n.Namespace)
	}

	if n.Code != "" {
		enc.AddString("code", n.Code)
	}

	if n.Message != "" {
		enc.AddString("message", n.Message)
	}

	if n.Caller != "" {
		enc.AddString("caller", n.Caller)
	}

	if n.Type != "" {
		enc.AddString("type", n.Type)
	}

	if len(n.Fields) > 0 {
		enc.AddObject("fields", zapcore.ObjectMarshalerFunc(func(obj zapcore.ObjectEncoder) error {
			for k, v := range n.Fields {
				obj.AddReflected(k, v)
			}
			return nil
		}))
	}

	return nil
}

func (e *errorStruct) getChain() []node {
	var chain []node
	var err error = e

	for err != nil {
		switch v := err.(type) {

		case *errorStruct:
			chain = append(chain, node{
				Namespace: v.Namespace,
				Code:      v.Code,
				Message:   v.Message,
				Fields:    v.Fields,
				Caller:    v.Caller,
				Type:      v.Type,
			})
		default:
			chain = append(chain, node{
				Type:    reflect.TypeOf(err).String(),
				Message: err.Error(),
			})
		}

		err = errors.Unwrap(err)
	}

	return chain
}

type Errorx interface {
	New(code, msg string, params ...any) error
	Wrap(err error, code, message string, params ...any) error
}

type errorx struct {
	namespace string
}

func (o *errorx) New(code, msg string, params ...any) error {
	fields := make(map[string]any)

	if len(params)%2 != 0 {
		panic("params must be key-value pairs")
	}

	for i := 0; i < len(params); i += 2 {
		key, ok := params[i].(string)
		if !ok {
			panic(fmt.Sprintf("param key at index %d is not string: %T", i, params[i]))
		}

		fields[key] = params[i+1]
	}

	return &errorStruct{
		Namespace: o.namespace,
		Code:      code,
		Message:   msg,
		Fields:    fields,
		Caller:    caller(2),
	}
}

func (o *errorx) Wrap(err error, code, message string, params ...any) error {
	errStruct := o.New(code, message, params...).(*errorStruct)
	errStruct.Err = err
	errStruct.Caller = caller(2)
	return errStruct
}

func NewNamespace(namespace string) Errorx {
	return &errorx{
		namespace: namespace,
	}
}
