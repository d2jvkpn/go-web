package api

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/d2jvkpn/go-web/pkg/misc"
	. "github.com/d2jvkpn/go-web/pkg/resp"

	"github.com/gin-gonic/gin"
)

func HandleUploadFile(ctx *gin.Context) {
	var (
		err        error
		fn, target string
		user, dir  string
		fileHeader *multipart.FileHeader
		files      []*multipart.FileHeader
		form       *multipart.Form
		now        time.Time
	)

	// fileHeader, err = ctx.FormFile("file")
	if form, err = ctx.MultipartForm(); err != nil {
		JSON(ctx, nil, ErrParseFailed(err))
		return
	}
	files = form.File["files"]
	user = ctx.GetString("User")
	now = time.Now()

	for _, fileHeader = range files {
		fn = fileHeader.Filename
		if fileHeader.Size == 0 {
			msg := "file is empty: " + fn
			JSON(ctx, nil, ErrInvalidParameter(fmt.Errorf(msg), msg))
			return
		}

		if !_FilenameRE.Match([]byte(fn)) {
			msg := "invalid filename: " + fn
			JSON(ctx, nil, ErrInvalidParameter(fmt.Errorf(msg), msg))
			return
		}

		if fileHeader.Size > HTTP_MaxFileHeaderSize {
			msg := "file size is too large: " + fn
			JSON(ctx, nil, ErrBadRequest(fmt.Errorf(msg), Msg(msg)))
			return
		}
	}

	dir = filepath.Join("data", "uploads", user, now.Format("2006-01-02"))
	if err = os.MkdirAll(dir, 0755); err != nil {
		JSON(ctx, nil, ErrServerError(err))
		return
	}

	for _, fileHeader = range files {
		fn = fileHeader.Filename
		ext := filepath.Ext(fn)
		fn = fileHeader.Filename[0 : len(fn)-len(ext)]

		target = filepath.Join(dir, fmt.Sprintf(
			"%s.%d_%s%s", fn, now.UnixMilli(), misc.RandString(16), ext,
		))

		log.Printf(
			"receiving file: source=%q, size=%s\n",
			fileHeader.Filename, misc.FileSize2Str(fileHeader.Size),
		)

		if err = ctx.SaveUploadedFile(fileHeader, target); err != nil {
			log.Printf("save file error: %q, %v\n", target, err)
			JSON(ctx, nil, ErrServerError(fmt.Errorf("save %s: %w", fn, err)))
			return
		}
	}

	JSON(ctx, gin.H{"url": strings.Replace(target, "data/", "/", 1)}, nil)
}
