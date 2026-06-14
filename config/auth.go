package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultResendFrom = "FreeStyleTarot <verify@freestyletarot.fun>"

type AuthConfig struct {
	DatabaseURL        string
	DatabaseMigrateURL string // 可选；PaaS 仅 IPv4 时可与 DATABASE_URL 同为 Session pooler（5432）
	JWTSecret          string
	JWTAccessTTL       time.Duration
	ResendAPIKey       string
	ResendFrom         string
	CORSOrigin         string
	ForceLogin         bool // 来自 config.yaml
	SkipVerify         bool // 来自 config.yaml
}

var Auth AuthConfig

func loadAuthConfig() {
	Auth.DatabaseURL = os.Getenv("DATABASE_URL")
	Auth.DatabaseMigrateURL = os.Getenv("DATABASE_MIGRATE_URL")
	if Auth.DatabaseMigrateURL == "" {
		Auth.DatabaseMigrateURL = Auth.DatabaseURL
	}
	Auth.JWTSecret = os.Getenv("JWT_SECRET")
	Auth.ResendAPIKey = os.Getenv("RESEND_API_KEY")
	Auth.CORSOrigin = os.Getenv("CORS_ORIGIN")
	if Auth.CORSOrigin == "" {
		Auth.CORSOrigin = "http://localhost:5173"
	}

	Auth.ForceLogin = GlobalConfig.Auth.ForceLogin
	Auth.SkipVerify = GlobalConfig.Auth.SkipVerify

	Auth.ResendFrom = normalizeResendFrom(os.Getenv("RESEND_FROM"))
	validateResendConfig()

	ttlStr := os.Getenv("JWT_ACCESS_TTL")
	if ttlStr == "" {
		Auth.JWTAccessTTL = 7 * 24 * time.Hour
	} else {
		d, err := time.ParseDuration(ttlStr)
		if err != nil {
			days, parseErr := strconv.Atoi(ttlStr)
			if parseErr != nil {
				Auth.JWTAccessTTL = 7 * 24 * time.Hour
			} else {
				Auth.JWTAccessTTL = time.Duration(days) * 24 * time.Hour
			}
		} else {
			Auth.JWTAccessTTL = d
		}
	}
}

func normalizeResendFrom(raw string) string {
	s := strings.TrimSpace(raw)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			s = strings.TrimSpace(s[1 : len(s)-1])
		}
	}
	if s == "" {
		return defaultResendFrom
	}
	return s
}

func validateResendConfig() {
	if Auth.SkipVerify {
		return
	}
	if Auth.ResendAPIKey == "" {
		log.Fatal("RESEND_API_KEY 未设置（skip_verify=false 时必需）")
	}
	if Auth.ResendFrom == "" {
		log.Fatal("RESEND_FROM 未设置（skip_verify=false 时必需）")
	}
	if strings.Contains(strings.ToLower(Auth.ResendFrom), "@resend.dev") {
		log.Fatal("RESEND_FROM 不能使用 resend.dev 测试域名，请使用已在 Resend 验证的自定义域名")
	}
	log.Printf("Resend 发件人: %s", Auth.ResendFrom)
}
