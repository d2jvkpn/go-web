package aliyun

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	// "time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssClient struct {
	*oss.Client
	config Config
}

func (client *OssClient) UploadLocal(fp, subpath string, options ...oss.Option) (
	link string, err error) {
	var (
		bucket *oss.Bucket
		config Config
	)

	config = client.config
	if bucket, err = client.Bucket(config.Bucket); err != nil {
		return "", err
	}

	// urlpath = strings.Trim(fmt.Sprintf("%s/%s", strings.Trim(config.Path, "/"), subpath), "/")
	// additional slash in middle will cause error: The specified object is not valid.
	// subpath with slash tail will create a new directory
	if subpath, err = ValidSubpath(subpath); err != nil {
		return "", err
	}

	if err = bucket.PutObjectFromFile(subpath, fp, options...); err != nil {
		return "", err
	}

	return client.config.Url(subpath), nil
}

func (client *OssClient) UploadFromUrl(urlSrc string, subpath string, options ...oss.Option) (
	link string, fsize int64, err error) {
	var (
		httpRes *http.Response
		bucket  *oss.Bucket
		config  Config
		ossReq  *oss.PutObjectRequest
		// ossRes  *oss.Response
	)

	config = client.config
	if bucket, err = client.Bucket(config.Bucket); err != nil {
		return "", 0, err
	}

	if subpath, err = ValidSubpath(subpath); err != nil {
		return "", 0, err
	}

	if httpRes, err = http.Head(urlSrc); err != nil {
		return "", 0, err
	}
	fsize, _ = strconv.ParseInt(httpRes.Header.Get("Content-Length"), 10, 64)

	if httpRes, err = http.Get(urlSrc); err != nil {
		return "", 0, err
	}
	// ?? no Content-Length here
	// fsize, _ = strconv.ParseInt(httpRes.Header.Get("Content-Length"), 10, 64)
	if httpRes.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("http status: %d, %s", httpRes.StatusCode, httpRes.Status)
	}
	defer httpRes.Body.Close()

	ossReq = &oss.PutObjectRequest{
		ObjectKey: subpath,
		Reader:    httpRes.Body,
	}

	if _, err = bucket.DoPutObject(ossReq, options); err != nil {
		return "", 0, err
	}
	// fmt.Printf("%#v\n", ossRes)

	return client.config.Url(subpath), fsize, nil
}

func (client *OssClient) PutObject(reader io.Reader, subpath string, options ...oss.Option) (
	link string, err error) {
	var (
		bucket *oss.Bucket
		config Config
		ossReq *oss.PutObjectRequest
		// ossRes  *oss.Response
	)

	if subpath, err = ValidSubpath(subpath); err != nil {
		return "", err
	}

	config = client.config
	if bucket, err = client.Bucket(config.Bucket); err != nil {
		return "", err
	}

	ossReq = &oss.PutObjectRequest{
		ObjectKey: subpath,
		Reader:    reader,
	}

	if _, err = bucket.DoPutObject(ossReq, options); err != nil {
		return "", err
	}
	// fmt.Printf("%#v\n", ossRes)
	return client.config.Url(subpath), nil
}

func (client *OssClient) Upload(fp string, subpath string, overWrite bool) (
	link string, fsize int64, err error) {

	var (
		file *os.File
		fi   fs.FileInfo
	)

	if strings.HasPrefix(fp, "http://") || strings.HasPrefix(fp, "https://") {
		return client.UploadFromUrl(fp, subpath, oss.ForbidOverWrite(!overWrite))
	}

	////
	if file, err = os.Open(fp); err != nil {
		return "", 0, err
	}

	if fi, err = file.Stat(); err != nil {
		file.Close()
		return "", 0, err
	}
	fsize = fi.Size()
	file.Close()

	if link, err = client.UploadLocal(fp, subpath, oss.ForbidOverWrite(!overWrite)); err != nil {
		return "", 0, err
	}

	return link, fsize, nil
}

func (client *OssClient) UploadDir(d string, dir string) (
	link string, err error) {

	//	if d, err = filepath.Abs(d); err != nil {
	//		return "", err
	//	}

	err = filepath.Walk(d, func(fp string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// fmt.Println(fp, info.Size())
		if info.IsDir() {
			return nil
		}
		// fmt.Println("~~~", time.Now().Format(time.RFC3339), fp)
		t := filepath.Join(dir, strings.Replace(fp, d, "", 1))
		_, err = client.UploadLocal(fp, t)
		return err
	})

	if err != nil {
		return "", err
	}

	// link is a dir path, try to access html like dir + "index/html"
	return client.config.Url(strings.Trim(dir, "/")), err
}

func (client *OssClient) UploadDirV2(d string, dir string, conc uint) (
	link string, err error) {

	//	if d, err = filepath.Abs(d); err != nil {
	//		return "", err
	//	}

	if conc == 0 {
		conc = 1
	}
	n, done, once := uint(0), make(chan bool, conc), new(sync.Once)

	// err = filepath.Walk(...)
	filepath.Walk(d, func(fp string, info fs.FileInfo, e1 error) error {
		if e1 != nil {
			return e1
		}

		// fmt.Println(fp, info.Size())
		if err != nil || info.IsDir() {
			return nil
		}

		if n++; n > conc {
			<-done
			n--
		}

		go func() {
			t := filepath.Join(dir, strings.Replace(fp, d, "", 1))
			// fmt.Println("~~~", time.Now().Format(time.RFC3339), fp)
			if _, e2 := client.UploadLocal(fp, t); e2 != nil {
				once.Do(func() {
					err = e2
				})
			}
			done <- true
		}()
		return nil
	})

	for ; n > 0; n-- {
		<-done
	}

	if err != nil {
		return "", err
	}

	// link is a dir path, try to access html like dir + "index/html"
	return client.config.Url(strings.Trim(dir, "/")), err
}

func (client *OssClient) CopyFile(src, target string, options ...oss.Option) (
	code string, err error) {
	var (
		bucket       *oss.Bucket
		config       Config
		serviceError oss.ServiceError
	)

	config = client.config
	if bucket, err = client.Bucket(config.Bucket); err != nil {
		serviceError, _ = err.(oss.ServiceError)
		return serviceError.Code, err
	}

	// result oss.CopyObjectResult
	if _, err = bucket.CopyObject(src, target, options...); err != nil {
		serviceError, _ = err.(oss.ServiceError)
		return serviceError.Code, err
	}

	return "", nil
}
