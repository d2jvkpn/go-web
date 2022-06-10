package resp

import (
	"net/http"
)

/// code ranges:
// ...=-100 has no right
// -99=-1   invalid request
// 0        ok
// 1..=100  business error
// 101...   unexpected error

func ErrBadRequest(cause error, opts ...Option) (err *HttpError) {
	opts = append(opts, Skip(2))
	return NewHttpError(cause, http.StatusBadRequest, -1, opts...)
}
