package api

import (
	"time"

	. "github.com/d2jvkpn/go-web/pkg/resp"

	"github.com/gin-gonic/gin"
)

func Load(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	open := rg.Group("/api/open", handlers...)

	open.GET("/timeout", func(ctx *gin.Context) {
		time.Sleep(20 * time.Second)
		JSON(ctx, gin.H{"code": 0, "msg": "ok"}, nil)
	})

	open.GET("/panic", func(ctx *gin.Context) {
		a, b := 1, 0
		result := a / b
		JSON(ctx, gin.H{"result": result}, nil)
	})

	open.POST("/login", login)
	open.GET("/hello", hello)
	open.GET("/hello/:name", hello)

	open.POST("/register", func(ctx *gin.Context) {
		data := struct {
			User     string `json:"user"`
			Password string `json:"password"`
		}{}

		var (
			err  error
			hErr *HttpError
		)

		if err = ctx.BindJSON(&data); err != nil {
			JSON(ctx, nil, ErrParseFailed(err))
			return
		}

		if hErr = _UsersData.Register(data.User, data.Password); hErr != nil {
			JSON(ctx, nil, hErr)
			return
		}

		Ok(ctx)
	})

	loadAuthApi(rg, append(handlers, basicAuth)...)
}

func loadAuthApi(rg *gin.RouterGroup, handlers ...gin.HandlerFunc) {
	// auth := rg.Group("/basic_auth", gin.BasicAuth(ReadAccounts())) // clone map
	// auth := rg.Group("/auth", basic_auth)
	auth := rg.Group("/api/auth", handlers...)

	user := auth.Group("/user")
	user.PUT("/upload", HandleUploadFile)

	user.POST("/unregister", func(ctx *gin.Context) {
		user := ctx.GetString(KEY_User)
		_UsersData.Unregister(user)
		Ok(ctx)
	})
}
