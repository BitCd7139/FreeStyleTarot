package db

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var Redis *redis.Client

// InitRedis 从 REDIS_URL 初始化 Redis（required=false 时跳过）
func InitRedis(required bool) {
	if !required {
		zap.S().Infow("skip_verify 已开启，跳过 Redis 初始化")
		return
	}

	url := os.Getenv("REDIS_URL")
	if url == "" {
		zap.S().Fatalw("REDIS_URL 未设置（启用真实验证码时必需）")
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		zap.S().Fatalw("REDIS_URL 解析失败", "error", err)
	}

	Redis = redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Redis.Ping(ctx).Err(); err != nil {
		zap.S().Fatalw("Redis Ping 失败", "error", err)
	}

	zap.S().Infow("Redis 连接成功")
}
