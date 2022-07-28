package internal

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/d2jvkpn/go-web/internal/services/api"
	"github.com/d2jvkpn/go-web/internal/services/site"
	"github.com/d2jvkpn/go-web/internal/services/ws"
	"github.com/d2jvkpn/go-web/internal/settings"
	"github.com/d2jvkpn/go-web/pkg/misc"
	"github.com/d2jvkpn/go-web/pkg/resp"
	"github.com/d2jvkpn/go-web/pkg/wrap"

	"github.com/gin-gonic/gin"
)

var (
	//go:embed static/assets/favicon.ico
	favicon embed.FS
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
		// engi.Use(gin.Recovery())
	} else {
		engi = gin.Default()
	}
	engi.RedirectTrailingSlash = false
	engi.MaxMultipartMemory = HTTP_MaxMultipartMemory
	rg = &engi.RouterGroup

	// engi.LoadHTMLGlob("templates/*.tmpl")
	tmpl, err = template.ParseFS(_Templates, "templates/*.html", "templates/*/*.html")
	if err != nil {
		return nil, err
	}
	engi.SetHTMLTemplate(tmpl)
	engi.Use(wrap.Cors("*"))

	///
	engi.NoRoute(func(ctx *gin.Context) {
		// ctx.AbortWithStatus(http.StatusNotFound)
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": -1, "msg": "router not found", "data": gin.H{},
		})
	})

	rg.StaticFileFS("/favicon.ico", "static/assets/favicon.ico", http.FS(favicon))
	rg.GET("/healthy", wrap.Healthy)
	rg.GET("/nts", gin.WrapF(misc.NTSFunc(3)))

	aipHandlers := []gin.HandlerFunc{resp.NewLogHandler(settings.Logger, "api")}
	if p := _Config.GetString("prometheus.path"); p != "" { // /prometheus
		rg.GET(p, wrap.PrometheusFunc)
		aipHandlers = append(aipHandlers, wrap.NewPrometheusMonitor("api"))
	}

	wrap.Pprof(rg) // TODO: more middlewares

	///
	if fsys, err = fs.Sub(_Static, "static"); err != nil {
		return nil, err
	}
	static := engi.RouterGroup.Group("/static", wrap.CacheControl(3600))
	static.StaticFS("/", http.FS(fsys))
	// ?? w.Header().Set("Cache-Control", "public, max-age=3600")
	// bts, _ := _Static.ReadFile("static/favicon.png")
	// engi.RouterGroup.GET("/favicon.ico", "image/x-icon", "favicon.ico", wrap.ServeFile(bts))

	for i := range _StaticDirs {
		if err = _StaticDirs[i](rg); err != nil {
			return nil, err
		}
	}

	site.Load(rg)
	ws.Load(rg, wrap.WsUpgrade, wrap.NewPrometheusMonitor("ws"))
	api.Load(rg, aipHandlers...)

	return
}
