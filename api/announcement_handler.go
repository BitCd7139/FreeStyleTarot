package api

import (
	"net/http"

	"FreeStyleTarot/config"

	"github.com/gin-gonic/gin"
)

// HandleAnnouncement GET /announcement — 返回当前进程公告与 boot_id
func HandleAnnouncement(c *gin.Context) {
	ann := config.GlobalConfig.Announcement
	c.JSON(http.StatusOK, gin.H{
		"enabled":  ann.Enabled,
		"boot_id":  config.ServerBootID,
		"title":    ann.Title,
		"content":  ann.Content,
	})
}
