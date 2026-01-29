package errorx

import "errors"

func Wrap(err error, namespace, code, message string, params ...any) error {
	errStruct := NewNamespace(namespace).New(code, message, params...).(*errorStruct)
	errStruct.Err = err
	errStruct.Caller = caller(2)
	return errStruct
}

func IsErrx(err error) bool {
	var ae *errorStruct
	return errors.As(err, &ae)
}

func IsCode(err error, code string) bool {
	var ae *errorStruct
	if !errors.As(err, &ae) {
		return false
	}

	for err != nil {
		if e, ok := err.(*errorStruct); ok && e.Code == code {
			return true
		}
		err = errors.Unwrap(err)
	}

	return false
}

func IsNamespace(err error, ns string) bool {
	for err != nil {
		if e, ok := err.(*errorStruct); ok && e.Namespace == ns {
			return true
		}
		err = errors.Unwrap(err)
	}

	return false
}

func IsCodeNS(err error, namespace, code string) bool {
	for err != nil {
		if e, ok := err.(*errorStruct); ok {
			if e.Code == code && e.Namespace == namespace {
				return true
			}
		}
		err = errors.Unwrap(err)
	}
	return false
}

func As[T any](err error) (*T, bool) {
	var target *T
	if errors.As(err, &target) {
		return target, true
	}
	return nil, false
}
