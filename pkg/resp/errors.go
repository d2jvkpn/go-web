package resp

import (
	"net/http"
)

/// code ranges:
// ...=-100 has no right
// -99=-1   invalid request
// 0        ok
// 1..=99   business error
// 100...   unexpected error

func ErrBadRequest(cause error, opts ...Option) (err *HttpError) {
	opts = append(opts, Skip(2))
	err = NewHttpError(cause, http.StatusBadRequest, -1, opts...)
	if err.Msg == "" {
		err.Msg = "bad request"
	}
	return
}

func ErrParseFailed(cause error) (err *HttpError) {
	err = NewHttpError(cause, http.StatusBadRequest, -2, Skip(2))
	if err.Msg == "" {
		err.Msg = "parse data failed"
	}
	return
}

func ErrInvalidParameter(err error, msg string) (out *HttpError) {
	return NewHttpError(err, http.StatusBadRequest, -3, Skip(2))
}

////
func ErrNotFound(err error) (out *HttpError) {
	return NewHttpError(err, http.StatusNotFound, -4, Skip(2), Msg("not found"))
}

////
func ErrServerError(cause error) (err *HttpError) {
	return NewHttpError(
		cause,
		http.StatusInternalServerError,
		100,
		Skip(2),
		Msg("internal server error"),
	)
}
