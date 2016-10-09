package route

import (
	"github.com/labstack/echo"
	echoMw "github.com/labstack/echo/middleware"
	"gofile-restapi/api"
	"github.com/flyer5200/gofile-restapi/handler"
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
		v1.GET("/files/:pv", api.GetFile())
		v1.POST("/files/:pv", api.PostFile())
		v1.DELETE("/files/:pv", api.DeleteFiles())
	}
	return e
}
