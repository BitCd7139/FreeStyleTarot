package api

import (
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type testCustomAPIRequest struct {
	APIKey      string             `json:"api_key"`
	BaseURL     string             `json:"base_url"`
	Model       string             `json:"model"`
	Format      string             `json:"format"`
	StageParams *request.StageParams `json:"stage_params,omitempty"`
}

// HandleTestCustomAPI POST /api/custom-api/test
// Tests the user-provided API config by making a minimal LLM call.
func HandleTestCustomAPI(c *gin.Context) {
	var req testCustomAPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.APIKey == "" || req.BaseURL == "" || req.Model == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api_key, base_url, model 不能为空"})
		return
	}
	if req.Format == "" {
		req.Format = "openai"
	}
	if req.Format != "openai" && req.Format != "anthropic" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format 仅支持 openai 或 anthropic"})
		return
	}

	cfg := request.CustomAPIConfig{
		APIKey:      req.APIKey,
		BaseURL:     req.BaseURL,
		Model:       req.Model,
		Format:      req.Format,
		StageParams: req.StageParams,
	}

	reply, err := service.TestCustomAPI(c.Request.Context(), cfg)
	if err != nil {
		zap.S().Warnw("custom API test failed", "error", err, "format", req.Format, "base_url", req.BaseURL, "model", req.Model)
		c.JSON(http.StatusOK, gin.H{
			"ok":      false,
			"error":   err.Error(),
			"message": "测试失败：" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"message": "测试成功，API 可正常使用",
		"reply":   reply,
	})
}
