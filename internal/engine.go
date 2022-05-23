package internal

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/d2jvkpn/goapp/internal/api"
	"github.com/d2jvkpn/goapp/internal/site"
	"github.com/d2jvkpn/goapp/internal/ws"
	"github.com/d2jvkpn/goapp/pkg/misc"

	"github.com/gin-gonic/gin"
)

var (
	//go:embed static
	_Static embed.FS
	//go:embed templates
	_Templates embed.FS
)

func NewEngine(release bool) (engi *gin.Engine, err error) {
	var (
		tmpl *template.Template
		fsys fs.FS
	)

	//
	if release {
		gin.SetMode(gin.ReleaseMode)
		engi = gin.New()
		engi.Use(gin.Recovery())
	} else {
		engi = gin.Default()
	}
	engi.RedirectTrailingSlash = false

	// engi.LoadHTMLGlob("templates/*.tmpl")
	if tmpl, err = template.ParseFS(_Templates, "templates/*.html"); err != nil {
		return nil, err
	}
	engi.SetHTMLTemplate(tmpl)
	engi.Use(misc.Cors)

	engi.NoRoute(func(ctx *gin.Context) {
		// ctx.AbortWithStatus(http.StatusNotFound)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": -1, "msg": "router not found", "data": gin.H{},
		})
	})

	if fsys, err = fs.Sub(_Static, "static"); err != nil {
		return nil, err
	}
	engi.RouterGroup.StaticFS("/site/static", http.FS(fsys))
	// bts, _ := _Static.ReadFile("static/favicon.png")
	// engi.RouterGroup.GET("/favicon.ico", "image/x-icon", "favicon.ico", misc.ServeFile(bts))

	//
	rg := &engi.RouterGroup
	api.Load(rg)
	ws.Load(rg, misc.WsUpgrade)
	site.Load(rg)

	return
}
