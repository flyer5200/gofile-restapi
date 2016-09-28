package route

import (
	"github.com/labstack/echo"
	echoMw "github.com/labstack/echo/middleware"
	"gofile-restapi/api"
	"gofile-restapi/handler"
)

func Init() *echo.Echo {
	e := echo.New()
	e.Debug()

	// Set Bundle MiddleWare
	e.Use(echoMw.Logger())
	e.Use(echoMw.Gzip())
	e.Use(echoMw.CORSWithConfig(echoMw.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding},
	}))
	e.SetHTTPErrorHandler(handler.JSONHTTPErrorHandler)

	// Routes
	v1 := e.Group("/api/v1")
	{
		v1.GET("/files/:id/:pv", api.ListDirs())
		v1.POST("/files", api.PostFile())
		v1.GET("/file/:path", api.GetFile())
		v1.DELETE("/file/:path", api.DeleteFiles())
	}
	return e
}
