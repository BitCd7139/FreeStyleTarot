package middleware

import (
	"net/http"
	"time"

	"FreeStyleTarot/config"
	"FreeStyleTarot/service/auth"

	"github.com/gin-gonic/gin"
)

var quotaService = auth.NewService()

// PredictQuota 免费用户 1 小时内限 1 次占卜
func PredictQuota() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Auth.ForceLogin {
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
