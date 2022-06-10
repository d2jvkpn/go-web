package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	KeyError     = "error"
	KeyRequestId = "requestId"
	KeyUserId    = "userId"
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
	d2["data"], d2["code"], d2["msg"] = gin.H{}, err.Code, err.Msg
	ctx.JSON(err.HttpCode, d2)

	return
}

func Ok(ctx *gin.Context) {
	JSON(ctx, nil, nil)
}

func BadRequest(ctx *gin.Context, cause error, msgs ...string) {
	var opts []Option

	opts = make([]Option, 0, 2)
	if len(msgs) > 0 {
		opts = append(opts, Msg(msgs[0]))
	} else {
		opts = append(opts, Msg("bad request"))
	}
	opts = append(opts, Skip(2))

	JSON(ctx, nil, NewHttpError(cause, http.StatusBadRequest, -1, opts...))
}
