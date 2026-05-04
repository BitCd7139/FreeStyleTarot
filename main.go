package main

import (
	"FreeStyleTarot/api"
	"FreeStyleTarot/config"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	_ = godotenv.Load()
	config.InitConfig()

	if strings.ToLower(config.GlobalConfig.Server.Mode) == "release" {
		// --- 生产环境 ---
		gin.SetMode(gin.ReleaseMode)
		logger, _ := zap.NewProduction()
		zap.ReplaceGlobals(logger)

	} else {
		// --- 开发环境 ---
		gin.SetMode(gin.DebugMode)
		logconf := zap.NewDevelopmentConfig()
		logconf.Encoding = "console"
		logconf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		logger, _ := logconf.Build()
		zap.ReplaceGlobals(logger)
	}

	r := gin.Default()

	r.POST("/predict", api.HandlePredictStream)
	r.POST("/prompt", api.HandlePrompt)

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
