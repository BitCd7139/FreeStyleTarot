package middleware

import (
	"net/http"
	"strings"

	"FreeStyleTarot/config"
	"FreeStyleTarot/service/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const UserIDKey = "user_id"

var authService = auth.NewService()

// tryParseToken 有 Bearer token 时解析并注入 user_id，失败则忽略
func tryParseToken(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" || !strings.HasPrefix(header, "Bearer ") {
		return
	}
	tokenStr := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	claims, err := authService.ParseToken(tokenStr)
	if err != nil {
		return
	}
	c.Set(UserIDKey, claims.UserID)
}

// AuthRequired 解析 JWT；force_login=false 时允许匿名通过（仍尝试解析 token）
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Auth.ForceLogin {
			tryParseToken(c)
			c.Next()
			return
		}

		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		claims, err := authService.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录已失效"})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

// AuthStrictRequired 必须登录（用于 /auth/me 等）
func AuthStrictRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		tokenStr := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		claims, err := authService.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "登录已失效"})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}

// GetUserID 从上下文读取 user_id
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	val, ok := c.Get(UserIDKey)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}
