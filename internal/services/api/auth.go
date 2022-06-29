package api

import (
	"encoding/base64"
	"fmt"
	"strings"

	. "github.com/d2jvkpn/go-web/pkg/resp"

	"github.com/gin-gonic/gin"
)

var basicAuth gin.HandlerFunc = func(ctx *gin.Context) {
	authorization := ctx.GetHeader("Authorization")
	var (
		found    bool
		bts      []byte
		user     string
		password string
		err      error
		herr     *HttpError
	)

	if !strings.HasPrefix(authorization, "Basic ") {
		ctx.Header("Www-Authenticate", `Basic realm="Authorization Required"`)
		msg := "invalid header value"
		JSON(ctx, nil, ErrUnauthorized(fmt.Errorf(msg), msg))
		ctx.Abort()
		return
	}

	if bts, err = base64.StdEncoding.DecodeString(authorization[6:]); err != nil {
		msg := "invalid header value"
		JSON(ctx, nil, ErrUnauthorized(err, msg))
		ctx.Abort()
		return
	}

	if user, password, found = strings.Cut(string(bts), ":"); !found {
		msg := "invalid header value"
		JSON(ctx, nil, ErrUnauthorized(fmt.Errorf(msg), msg))
		ctx.Abort()
		return
	}

	if herr = _UsersData.Verify(user, password); herr != nil {
		JSON(ctx, nil, herr)
		ctx.Abort()
		return
	}

	ctx.Set(KEY_User, user)
	ctx.Next()
}
