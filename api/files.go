package api

import (
	"github.com/labstack/echo"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var result struct {
	status int
	msg    string
}

var basePath = "/var/lib/kubelet/pods/"
var volumePath = "/volumes/kubernetes.io~glusterfs/"

// 获取文件流
func GetFile() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// 通过名称
		path := c.Param("path")
		file, error := ioutil.ReadFile(path)
		if error != nil {
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
		files, error := ioutil.ReadDir(basePath + id + volumePath + pv)
		if error != nil && files != nil {
			return c.JSON(fasthttp.StatusOK, files)
		}
		return c.JSON(fasthttp.StatusOK, nil)
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
