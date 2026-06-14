package db

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

// NormalizeDatabaseURL 修正常见粘贴错误，并为 Supabase Session pooler 补齐默认参数。
func NormalizeDatabaseURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("DATABASE_URL 为空")
	}

	raw = fixCommonSchemeTypos(raw)

	if !strings.Contains(raw, "://") {
		return "", fmt.Errorf("DATABASE_URL 缺少协议头，应以 postgresql:// 开头")
	}

	if strings.HasPrefix(raw, "postgres://") {
		raw = "postgresql" + strings.TrimPrefix(raw, "postgres")
	}

	raw = repairPoolerURLInPath(raw)

	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("DATABASE_URL 解析失败: %w", err)
	}

	if u.Scheme != "postgresql" {
		return "", fmt.Errorf("DATABASE_URL 协议须为 postgresql://，当前为 %q", u.Scheme)
	}

	if !IsSupabaseSessionPooler(u) && strings.Contains(u.Path, ".pooler.supabase.com") {
		return "", fmt.Errorf("DATABASE_URL 格式错误：pooler 地址被放在了路径里，请只保留一条 Session pooler URI（用户名 postgres.[ref]@*.pooler.supabase.com:5432）")
	}

	if IsSupabaseSessionPooler(u) {
		if err := validateSupabaseSessionPooler(u); err != nil {
			return "", err
		}
		ensureQueryDefault(u, "sslmode", "require")
	}

	return u.String(), nil
}

// ToMigrateURL 将应用 DSN 转为 golang-migrate pgx/v5 驱动格式。
func ToMigrateURL(dsn string) string {
	switch {
	case strings.HasPrefix(dsn, "postgresql://"):
		return "pgx5://" + strings.TrimPrefix(dsn, "postgresql://")
	case strings.HasPrefix(dsn, "postgres://"):
		return "pgx5://" + strings.TrimPrefix(dsn, "postgres://")
	default:
		return dsn
	}
}

// IsSupabaseSessionPooler 判断是否为 Supabase Shared Pooler（Session 模式，IPv4）。
func IsSupabaseSessionPooler(u *url.URL) bool {
	host := u.Hostname()
	return strings.Contains(host, ".pooler.supabase.com")
}

func fixCommonSchemeTypos(raw string) string {
	switch {
	case strings.HasPrefix(raw, "database:/"):
		return "postgresql://" + strings.TrimPrefix(raw, "database:/")
	case strings.HasPrefix(raw, "database://"):
		return "postgresql://" + strings.TrimPrefix(raw, "database://")
	default:
		return raw
	}
}

// repairPoolerURLInPath 修复两条 URI 拼在一起时 pooler 落在 path 里的情况。
func repairPoolerURLInPath(raw string) string {
	const marker = ".pooler.supabase.com"
	idx := strings.Index(raw, marker)
	if idx == -1 {
		return raw
	}

	schemeEnd := strings.Index(raw, "://")
	if schemeEnd == -1 {
		return raw
	}

	afterScheme := raw[schemeEnd+3:]
	firstAt := strings.Index(afterScheme, "@")
	if firstAt == -1 {
		return raw
	}

	hostPart := afterScheme[firstAt+1:]
	hostEnd := strings.IndexAny(hostPart, "/?#")
	if hostEnd == -1 {
		hostEnd = len(hostPart)
	}
	if strings.Contains(hostPart[:hostEnd], marker) {
		return raw
	}

	start := strings.LastIndex(raw[:idx], "postgres.")
	if start == -1 {
		return raw
	}

	endRel := strings.Index(raw[idx:], "/postgres")
	if endRel == -1 {
		return raw
	}
	end := idx + endRel + len("/postgres")

	return raw[:schemeEnd] + "://" + raw[start:end]
}

func validateSupabaseSessionPooler(u *url.URL) error {
	host := u.Hostname()
	port := u.Port()
	if port == "" {
		port = "5432"
		u.Host = net.JoinHostPort(host, port)
	}

	switch port {
	case "5432":
	case "6543":
		return fmt.Errorf("端口 6543 为 Transaction pooler；Session pooler（IPv4 PaaS）须使用 5432")
	default:
		return fmt.Errorf("Supabase Session pooler 端口应为 5432，当前为 %s", port)
	}

	user := ""
	if u.User != nil {
		user = u.User.Username()
	}
	if !strings.HasPrefix(user, "postgres.") {
		return fmt.Errorf("Session pooler 用户名须为 postgres.[project-ref]，当前为 %q", user)
	}

	return nil
}

func ensureQueryDefault(u *url.URL, key, value string) {
	q := u.Query()
	if q.Get(key) == "" {
		q.Set(key, value)
		u.RawQuery = q.Encode()
	}
}
