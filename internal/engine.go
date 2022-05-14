package internal

import (
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
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
		ctx.JSON(http.StatusNotFound, gin.H{"code": -1, "message": "not found", "data": nil})
	})

	if fsys, err = fs.Sub(_Static, "static"); err != nil {
		return nil, err
	}
	engi.RouterGroup.StaticFS("/static", http.FS(fsys))

	//
	LoadAPI(engi)

	return
}
