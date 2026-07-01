package db

import (
	"context"
	"database/sql"
	"net/url"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

var Pool *sql.DB

// Init 从 DATABASE_URL 初始化连接池，失败则退出进程
func Init() {
	raw := os.Getenv("DATABASE_URL")
	if raw == "" {
		zap.S().Fatalw("DATABASE_URL 未设置")
	}

	dsn, err := NormalizeDatabaseURL(raw)
	if err != nil {
		zap.S().Fatalw("DATABASE_URL 无效", "error", err)
	}

	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		zap.S().Fatalw("数据库连接失败", "error", err)
	}

	applyPoolSettings(conn, dsn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		zap.S().Fatalw("数据库 Ping 失败", "error", err)
	}

	Pool = conn
	zap.S().Infow("数据库连接成功")
}

func applyPoolSettings(conn *sql.DB, dsn string) {
	maxOpen := 10
	maxIdle := 5
	lifetime := 30 * time.Minute

	if u, err := url.Parse(dsn); err == nil && IsSupabaseSessionPooler(u) {
		// Session pooler 侧已有连接复用，应用侧池子宜小，避免占满 Supavisor 配额
		maxOpen = 5
		maxIdle = 2
		lifetime = 10 * time.Minute
	}

	conn.SetMaxOpenConns(maxOpen)
	conn.SetMaxIdleConns(maxIdle)
	conn.SetConnMaxLifetime(lifetime)
}
