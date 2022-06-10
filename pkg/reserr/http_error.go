package reserr

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"
)

type HttpError struct {
	Cause    error  `json:"cause"`
	HttpCode int    `json:"httpCode"`
	Code     int    `json:"code"`
	Msg      string `json:"msg"`
	skip     int
}

type Option func(*HttpError)

func Msg(msg string) Option {
	return func(e *HttpError) {
		e.Msg = msg
	}
}

func Skip(skip int) Option {
	return func(e *HttpError) {
		e.skip = skip
	}
}

func ErrIsNil() (out *HttpError) {
	err := fmt.Errorf("not a HttpError")
	return NewHttpError(err, http.StatusInternalServerError, 100, Msg("not a HttpError"))
}

func NewHttpError(cause error, httpCode, code int, opts ...Option) (err *HttpError) {
	if cause == nil { // avoid panic
		return ErrIsNil()
	}

	// httpCode != http.StatusOK, code != 0
	err = &HttpError{HttpCode: httpCode, Code: code, Msg: "", skip: 1}
	for _, v := range opts {
		v(err)
	}

	fn, file, line, _ := runtime.Caller(err.skip)

	err.Cause = fmt.Errorf(
		"%s(%s:%d): %w", runtime.FuncForPC(fn).Name(),
		filepath.Base(file), line, cause,
	)

	return err
}
