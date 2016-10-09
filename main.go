package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/flyer5200/gofile-restapi/route"
	"github.com/labstack/echo/engine/fasthttp"
	"fmt"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

}

func main() {
	router := route.Init()
	fmt.Println("K8sVolume API Started! Listen on 0.0.0.0:8888")
	router.Run(fasthttp.New(":8888"))
}
