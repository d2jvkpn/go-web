package internal

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/d2jvkpn/go-web/internal/services/api"
	"github.com/d2jvkpn/go-web/internal/services/site"
	"github.com/d2jvkpn/go-web/internal/services/ws"
	"github.com/d2jvkpn/go-web/pkg/misc"
	"github.com/d2jvkpn/go-web/pkg/resp"

	"github.com/gin-gonic/gin"
)

func NewEngine(release bool) (engi *gin.Engine, err error) {
	var (
		tmpl *template.Template
		fsys fs.FS
		rg   *gin.RouterGroup
	)

	///
	if release {
		gin.SetMode(gin.ReleaseMode)
		engi = gin.New()
		engi.Use(gin.Recovery())
	} else {
		engi = gin.Default()
	}
	engi.RedirectTrailingSlash = false
	rg = &engi.RouterGroup

	// engi.LoadHTMLGlob("templates/*.tmpl")
	tmpl, err = template.ParseFS(_Templates, "templates/*.html", "templates/*/*.html")
	if err != nil {
		return nil, err
	}
	engi.SetHTMLTemplate(tmpl)
	engi.Use(misc.Cors("*"))

	///
	engi.NoRoute(func(ctx *gin.Context) {
		// ctx.AbortWithStatus(http.StatusNotFound)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": -1, "msg": "router not found", "data": gin.H{},
		})
	})

	rg.GET("/misc/healthy", misc.Healthy)
	rg.GET("/misc/nts", gin.WrapF(misc.NTSFunc(3)))
	misc.Pprof(rg) // TODO: more middlewares

	///
	if fsys, err = fs.Sub(_Static, "static"); err != nil {
		return nil, err
	}
	static := engi.RouterGroup.Group("/static", misc.CacheControl(3600))
	static.StaticFS("/", http.FS(fsys))
	// ?? w.Header().Set("Cache-Control", "public, max-age=3600")
	// bts, _ := _Static.ReadFile("static/favicon.png")
	// engi.RouterGroup.GET("/favicon.ico", "image/x-icon", "favicon.ico", misc.ServeFile(bts))
	site.Load(rg)
	ws.Load(rg, misc.WsUpgrade)

	api.Load(rg, resp.NewLogHandler(_ApiLogger), misc.NewPrometheusMonitor())

	return
}
