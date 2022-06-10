package reserr

import (
	"net/http"
)

func ErrBadRequest(cause error, opts ...Option) (err *HttpError) {
	return NewHttpError(cause, http.StatusBadRequest, -1, opts...)
}
