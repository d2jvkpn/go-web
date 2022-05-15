package internal

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/d2jvkpn/goapp/internal/api"

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

	// engi.LoadHTMLGlob("templates/*.tmpl")
	if tmpl, err = template.ParseFS(_Templates, "templates/*.tmpl"); err != nil {
		return nil, err
	}
	engi.SetHTMLTemplate(tmpl)
	engi.Use(Cors)

	engi.NoRoute(func(ctx *gin.Context) {
		// ctx.AbortWithStatus(http.StatusNotFound)
		ctx.JSON(http.StatusNotFound, gin.H{"code": -1, "msg": "not found", "data": nil})
	})

	if fsys, err = fs.Sub(_Static, "static"); err != nil {
		return nil, err
	}
	engi.RouterGroup.StaticFS("/static", http.FS(fsys))

	//
	api.LoadAPI(&engi.RouterGroup)

	return
}
