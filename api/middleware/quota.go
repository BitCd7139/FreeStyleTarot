package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"FreeStyleTarot/config"
	"FreeStyleTarot/service/auth"

	"github.com/gin-gonic/gin"
)

var quotaService = auth.NewService()

// requestHasCustomAPI peeks at the JSON body to detect a custom_api.api_key
// field without consuming the body for downstream handlers.
func requestHasCustomAPI(c *gin.Context) bool {
	if c.Request.Body == nil {
		return false
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return false
	}
	// Restore the body so handlers can still read it.
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var peek struct {
		CustomAPI *struct {
			APIKey string `json:"api_key"`
		} `json:"custom_api"`
	}
	if err := json.Unmarshal(bodyBytes, &peek); err != nil {
		return false
	}
	return peek.CustomAPI != nil && peek.CustomAPI.APIKey != ""
}

// PredictQuota 免费用户 1 小时内限 1 次占卜
func PredictQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Auth.ForceLogin {
			c.Next()
			return
		}

		// 自定义 API 模式：用户自带 token，跳过服务器配额限制
		if requestHasCustomAPI(c) {
			c.Next()
			return
		}

		userID, ok := GetUserID(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			return
		}

		allowed, waitMinutes, err := quotaService.CheckPredictAllowed(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !allowed {
			var nextAt *time.Time
			if u, err := quotaService.GetMe(c.Request.Context(), userID); err == nil {
				nextAt = u.NextPredictAt
			}
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":           "提问次数已达上限",
				"wait_minutes":    waitMinutes,
				"next_predict_at": nextAt,
			})
			return
		}

		c.Next()
	}
}
