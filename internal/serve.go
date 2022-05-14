package internal

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Serve(addr string) (err error) {
	var engi *gin.Engine

	if engi, err = NewEngine(_Release); err != nil {
		return err
	}

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
