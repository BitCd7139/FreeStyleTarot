package main

import (
	"FreeStyleTarot/api"
	"FreeStyleTarot/config"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()
	config.InitConfig()

	if strings.ToLower(config.GlobalConfig.Server.Mode) == "release" {
		gin.SetMode(gin.ReleaseMode)

		logconf := zap.NewDevelopmentConfig()
		logconf.Encoding = "console"
		logger, _ := logconf.Build()
		zap.ReplaceGlobals(logger)

	} else {
		gin.SetMode(gin.DebugMode)
		logger, _ := zap.NewProduction()
		zap.ReplaceGlobals(logger)
	}

	r := gin.Default()

	r.POST("/predict", api.HandlePredictStream) // for test

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if port[0] != ':' {
		port = ":" + port
	}
	err := r.Run(port)

	if err != nil {
		panic("服务器启动失败: " + err.Error())
	}
}
