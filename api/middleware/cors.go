package middleware

import (
	"net/http"

	"FreeStyleTarot/config"

	"github.com/gin-gonic/gin"
)

// CORS 允许前端 Origin 与 Authorization header
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := config.Auth.CORSOrigin

		if origin == allowed || allowed == "*" {
			c.Header("Access-Control-Allow-Origin", origin)
			if allowed == "*" && origin == "" {
				c.Header("Access-Control-Allow-Origin", "*")
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
