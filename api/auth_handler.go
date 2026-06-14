package api

import (
	"net/http"
	"strings"

	"FreeStyleTarot/api/middleware"
	"FreeStyleTarot/config"
	"FreeStyleTarot/service/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var authSvc = auth.NewService()

type sendCodeRequest struct {
	Email string `json:"email"`
}

type verifyRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Nickname string `json:"nickname"`
}

type verifyCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type completeCodeSignupRequest struct {
	Email           string `json:"email"`
	Nickname        string `json:"nickname"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type patchMeRequest struct {
	Nickname string `json:"nickname"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Code     string `json:"code"`
}

type resetPasswordRequest struct {
	Email           string `json:"email"`
	Code            string `json:"code"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

// HandleSendCode POST /auth/send-code
func HandleSendCode(c *gin.Context) {
	var req sendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(req.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱不能为空"})
		return
	}

	err := authSvc.SendCode(c.Request.Context(), req.Email)
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "秒后再试") {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": msg})
			return
		}
		zap.S().Errorw("发送验证码失败", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "验证码已发送"})
}

// HandleVerify POST /auth/verify
func HandleVerify(c *gin.Context) {
	var req verifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(req.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱不能为空"})
		return
	}
	if !config.Auth.SkipVerify && len(strings.TrimSpace(req.Code)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码不能为空"})
		return
	}

	result, err := authSvc.Verify(c.Request.Context(), req.Email, req.Code, req.Nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": result.Token,
		"user":  result.User,
	})
}

// HandleVerifyCode POST /auth/verify-code — 仅校验验证码
func HandleVerifyCode(c *gin.Context) {
	var req verifyCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(req.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱不能为空"})
		return
	}
	if !config.Auth.SkipVerify && len(strings.TrimSpace(req.Code)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码不能为空"})
		return
	}

	result, err := authSvc.VerifyCodeOnly(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp := gin.H{"needs_setup": result.NeedsSetup}
	if !result.NeedsSetup {
		resp["token"] = result.Token
		resp["user"] = result.User
	}
	c.JSON(http.StatusOK, resp)
}

// HandleCompleteCodeSignup POST /auth/complete-code-signup — 验证码通过后补全资料
func HandleCompleteCodeSignup(c *gin.Context) {
	var req completeCodeSignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(req.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱不能为空"})
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "两次输入的密码不一致"})
		return
	}

	result, err := authSvc.CompleteCodeSignup(
		c.Request.Context(),
		req.Email,
		req.Nickname,
		req.Password,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": result.Token,
		"user":  result.User,
	})
}

// HandleLogin POST /auth/login
func HandleLogin(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(req.Email)) == 0 || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱和密码不能为空"})
		return
	}

	result, err := authSvc.LoginWithPassword(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": result.Token,
		"user":  result.User,
	})
}

// HandleRegister POST /auth/register
func HandleRegister(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(req.Email)) == 0 || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱和密码不能为空"})
		return
	}
	if !config.Auth.SkipVerify && len(strings.TrimSpace(req.Code)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码不能为空"})
		return
	}

	result, err := authSvc.Register(c.Request.Context(), req.Email, req.Password, req.Nickname, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": result.Token,
		"user":  result.User,
	})
}

// HandleResetPassword POST /auth/reset-password — 邮箱验证码重置密码
func HandleResetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(req.Email)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱不能为空"})
		return
	}
	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
		return
	}
	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "两次输入的密码不一致"})
		return
	}
	if !config.Auth.SkipVerify && len(strings.TrimSpace(req.Code)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码不能为空"})
		return
	}

	err := authSvc.ResetPassword(c.Request.Context(), req.Email, req.Code, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "密码已重置，请使用新密码登录"})
}

// HandleGetMe GET /auth/me
func HandleGetMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	result, err := authSvc.GetMe(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// HandlePatchMe PATCH /auth/me
func HandlePatchMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req patchMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(strings.TrimSpace(req.Nickname)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "昵称不能为空"})
		return
	}

	user, err := authSvc.UpdateNickname(c.Request.Context(), userID, strings.TrimSpace(req.Nickname))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// HandleAuthConfig GET /auth/config — 前端读取非敏感 auth 配置
func HandleAuthConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"force_login": config.Auth.ForceLogin,
		"skip_verify": config.Auth.SkipVerify,
	})
}
