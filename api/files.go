package api

import (
	"github.com/labstack/echo"
	"github.com/valyala/fasthttp"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"gofile-restapi/config"
)

var result struct {
	Status int
	Msg    string
}

type (
	fileInfo struct{
		Name string       // base name of the file
		Size int64        // length in bytes for regular files; system-dependent for others
		ModTime time.Time // modification time
		IsDir bool        // abbreviation for Mode().IsDir()
	}
)

func GetFile() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// 通过名称
		pv := c.Param("pv")
		path := c.QueryParam("path")
		if(path == ""){
			return rangeDirs(c, config.PvLink[pv])
		}
		isdir := isFileOrDir(config.PvLink[pv] +"/"+  path, true)
		if(isdir){
			return rangeDirs(c, config.PvLink[pv] +"/"+  path)
		}
		return download(c, config.PvLink[pv] +"/"+  path, path)
	}
}

func download(c echo.Context, path string, filename string) error  {
	file, err := ioutil.ReadFile(path)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment;filename="+filename)
	c.Response().WriteHeader(http.StatusOK)
	c.Response().Write(file)
	return err
}

func rangeDirs(c echo.Context, path string)  error {
	results, _ := ioutil.ReadDir(path)
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

// 判断是文件还是目录，根据decideDir为true表示判断是否为目录；否则判断是否为文件
func isFileOrDir(filename string, decideDir bool) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	isDir := fileInfo.IsDir()
	if decideDir {
		return isDir
	}
	return !isDir
}

func PostFile() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		pv := c.Param("pv")
		path := c.QueryParam("path")
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
		dst, err := os.Create(config.PvLink[pv]+"/"+ path+ "/"+ file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		// Copy
		if _, err = io.Copy(dst, src); err != nil {
			return err
		}
		result.Status = fasthttp.StatusOK
		result.Msg = "文件上传成功!"
		return c.JSON(fasthttp.StatusOK, result)
	}
}

func DeleteFiles() echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		// 通过名称
		pv := c.Param("pv")
		path := c.QueryParam("path")
		if path != "" {
			os.RemoveAll(config.PvLink[pv] +"/"+  path)
			result.Status = fasthttp.StatusOK
			result.Msg = "文件删除成功!"
			return c.JSON(fasthttp.StatusOK, result)
		}
		result.Status = fasthttp.StatusInternalServerError
		result.Msg = "文件删除失败!"
		return c.JSON(fasthttp.StatusInternalServerError, result)
	}
}
