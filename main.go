package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/flyer5200/gofile-restapi/route"
	"github.com/labstack/echo/engine/fasthttp"
	"log"
	"gofile-restapi/config"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	router := route.Init()
	log.Println("PV-FileManager Started! Listen on", config.Config["bind"])
	router.Run(fasthttp.New(config.Config["bind"]))
}
