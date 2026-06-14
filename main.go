package main

import (
	"FreeStyleTarot/api"
	"FreeStyleTarot/api/middleware"
	"FreeStyleTarot/config"
	"FreeStyleTarot/db"
	"FreeStyleTarot/service"
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
	config.InitBootID()
	service.InitAgentConfig()

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

	// 1. 初始化数据库、Redis 并执行 migration
	db.Init()
	db.InitRedis(!config.Auth.SkipVerify)
	db.RunMigrations(config.Auth.DatabaseMigrateURL)

	r := gin.Default()
	r.Use(middleware.CORS())

	// 公告（无需登录）
	r.GET("/announcement", api.HandleAnnouncement)

	// Auth 路由（无需登录）
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/config", api.HandleAuthConfig)
		authGroup.POST("/send-code", api.HandleSendCode)
		authGroup.POST("/verify", api.HandleVerify)
		authGroup.POST("/verify-code", api.HandleVerifyCode)
		authGroup.POST("/complete-code-signup", api.HandleCompleteCodeSignup)
		authGroup.POST("/login", api.HandleLogin)
		authGroup.POST("/register", api.HandleRegister)
		authGroup.POST("/reset-password", api.HandleResetPassword)
		authGroup.GET("/me", middleware.AuthStrictRequired(), api.HandleGetMe)
		authGroup.PATCH("/me", middleware.AuthStrictRequired(), api.HandlePatchMe)
		authGroup.GET("/predict-history", middleware.AuthStrictRequired(), api.HandleGetPredictHistory)
	}

	// 占卜路由（需登录 + 配额）
	protected := r.Group("")
	protected.Use(middleware.AuthRequired(), middleware.PredictQuota())
	{
		protected.POST("/predict", api.HandlePredictStream)
		protected.POST("/prompt", api.HandlePrompt)
	}

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
