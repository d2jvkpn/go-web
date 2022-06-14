package resp

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

type HttpError struct {
	Cause    string `json:"cause"`
	cause    error
	HttpCode int    `json:"httpCode"`
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	skip     int
}

type Option func(*HttpError) bool

func Msg(msg string) Option {
	return func(e *HttpError) bool {
		e.Msg = msg
		return true
	}
}

func Skip(skip int) Option {
	return func(e *HttpError) bool {
		if skip > 1 {
			e.skip = skip
			return true
		}
		return false
	}
}

func NewHttpError(cause error, httpCode, code int, opts ...Option) (err *HttpError) {
	if cause == nil { // avoid panic
		return nil
	}

	// httpCode != http.StatusOK, code != 0
	err = &HttpError{HttpCode: httpCode, Code: code, Msg: "", skip: 1}
	for _, v := range opts {
		_ = v(err)
	}

	fn, file, line, _ := runtime.Caller(err.skip)

	err.cause = fmt.Errorf(
		"%s(%s:%d): %w", runtime.FuncForPC(fn).Name(),
		filepath.Base(file), line, cause,
	)
	err.Cause = err.cause.Error()

	return err
}

func (err *HttpError) Unwrap() error {
	return errors.Unwrap(err.cause)
}
