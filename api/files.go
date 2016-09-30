package api

import (
	"github.com/labstack/echo"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var result struct {
	status int
	msg    string
}

type (
	fileInfo struct{
		Name string       // base name of the file
		Size int64        // length in bytes for regular files; system-dependent for others
		ModTime time.Time // modification time
		IsDir bool        // abbreviation for Mode().IsDir()
	}
)

//var basePath = "/var/lib/kubelet/pods/"
var basePath = "d://"
var volumePath = "/volumes/kubernetes.io~glusterfs/"

// 获取文件流
func GetFile() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// 通过名称
		path := c.Param("path")
		name := c.Param("name")
		file, e := ioutil.ReadFile(path + name)
		if e != nil {
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
			c.Response().WriteHeader(http.StatusOK)
			c.Response().Write(file)
			c.Response().(http.Flusher).Flush()
		}
		return nil
	}
}

func ListDirs() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// 通过名称
		id := c.Param("id")
		pv := c.Param("pv")
		results, _ := ioutil.ReadDir(basePath + id + volumePath + pv)
		var l [] *fileInfo
		for _, file := range results{
			fileinfo := &fileInfo{
				Name: file.Name(),
				IsDir:file.IsDir(),
				ModTime:file.ModTime(),
				Size:file.Size(),
			}
			l = append(l, fileinfo)
		}
		return c.JSON(fasthttp.StatusOK, l)
	}
}

func PostFile() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// Source
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// Destination
		dst, err := os.Create(basePath + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
		result.status = fasthttp.StatusOK
		result.msg = "文件上传成功!"
		return c.JSON(fasthttp.StatusOK, result)
	}
}

func DeleteFiles() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// 通过名称
		path := c.Param("path")
		if path != "" {
			err := os.RemoveAll(basePath + path)
			if err != nil {
				result.status = fasthttp.StatusOK
				result.msg = "文件上传成功!"
				return c.JSON(fasthttp.StatusOK, result)
			}
		}
		result.status = fasthttp.StatusInternalServerError
		result.msg = "文件上传失败!"
		return c.JSON(fasthttp.StatusInternalServerError, result)
	}
}
