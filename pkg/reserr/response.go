package reserr

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	KeyError     = "error"
	KeyRequestId = "requestId"
)

func JSON(ctx *gin.Context, data any, err *HttpError) {
	var (
		requestId string
		d2        map[string]interface{}
	)

	requestId = ctx.GetString(KeyRequestId)
	d2 = gin.H{"code": 0, "msg": "ok", "requestId": requestId}
	if err == nil {
		if data == nil {
			d2["data"] = gin.H{}
		}
		ctx.JSON(http.StatusOK, d2)
		return
	}

	ctx.Set(KeyError, err)
	d2["data"], d2["code"], d2["msg"] = nil, err.Code, err.Msg
	ctx.JSON(err.HttpCode, d2)

	return
}

func Ok(ctx *gin.Context) {
	JSON(ctx, nil, nil)
}

func ErrBadRequest(ctx *gin.Context, cause error, msgs ...string) {
	err := NewHttpError(cause, http.StatusBadRequest, -1)

	if len(msgs) > 0 {
		err.Msg = msgs[0]
	}
	JSON(ctx, nil, err)
}
