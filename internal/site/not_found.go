package site

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func notFound(ctx *gin.Context) {
	ctx.HTML(http.StatusNotFound, "page_not_found", nil)
}
