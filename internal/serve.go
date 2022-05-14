package internal

import (
	"html/template"
	"io/fs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Serve(addr string) (err error) {
	var (
		tmpl *template.Template
		fsys fs.FS
		engi *gin.Engine
	)

	if _Relase {
		gin.SetMode(gin.ReleaseMode)
		engi = gin.New()
		engi.Use(gin.Recovery())
	} else {
		engi = gin.Default()
	}

	// engi.LoadHTMLGlob("templates/*.tmpl")
	if tmpl, err = template.ParseFS(_Templates, "templates/*.tmpl"); err != nil {
		return err
	}
	engi.SetHTMLTemplate(tmpl)
	engi.Use(Cors)

	engi.NoRoute(func(ctx *gin.Context) {
		// ctx.AbortWithStatus(http.StatusNotFound)
		ctx.JSON(http.StatusNotFound, gin.H{"code": -1, "message": "not found", "data": nil})
	})

	if fsys, err = fs.Sub(_Static, "static"); err != nil {
		return err
	}
	engi.RouterGroup.StaticFS("/static", http.FS(fsys))

	LoadAPI(engi)

	_Server = &http.Server{
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		MaxHeaderBytes:    4 << 20,
		Addr:              addr,
		Handler:           engi,
	}

	if err = _Server.ListenAndServe(); err == http.ErrServerClosed {
		err = nil
	}

	return
}
