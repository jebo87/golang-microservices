package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jebo87/golang-microservices/src/api/log/logruslog"
	"github.com/jebo87/golang-microservices/src/api/log/zaplog"
	"github.com/sirupsen/logrus"
)

var (
	router *gin.Engine
)

func init() {
	router = gin.Default()
}

func StartApp() {
	logruslog.Info("logrus: About to map the URLs", "step:1", "status:pending")
	zaplog.Info("zap: About to map the URLs", zaplog.Field("step", "1"), zaplog.Field("status", "pending"))
	mapURLs()
	logrus.Info("URLs mapped succesfully", "step:2", "status:executed")

	if err := router.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
