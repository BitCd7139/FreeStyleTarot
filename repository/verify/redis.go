package verify

import (
	"context"
	"errors"
	"strconv"
	"time"

	"FreeStyleTarot/db"

	"github.com/redis/go-redis/v9"
)

const (
	verifyKeyPrefix   = "auth:verify:"
	cooldownKeyPrefix = "auth:cooldown:"
	codeTTL           = 30 * time.Minute
	cooldownTTL       = 60 * time.Second
)

var errNotFound = errors.New("验证码不存在或已过期")

type Store struct{}

func NewStore() *Store {
	return &Store{}
}

// Record 验证码记录
type Record struct {
	Hash     string
	Attempts int
}

func verifyKey(email string) string   { return verifyKeyPrefix + email }
func cooldownKey(email string) string { return cooldownKeyPrefix + email }

// CooldownRemain 返回重发冷却剩余秒数，0 表示可发送
func (s *Store) CooldownRemain(ctx context.Context, email string) (int, error) {
	ttl, err := db.Redis.TTL(ctx, cooldownKey(email)).Result()
	if err != nil {
		return 0, err
	}
	if ttl <= 0 {
		return 0, nil
	}
	secs := int(ttl.Seconds())
	if secs < 1 {
		secs = 1
	}
	return secs, nil
}

// SaveCode 写入验证码 hash，并设置 60s 重发冷却
func (s *Store) SaveCode(ctx context.Context, email, codeHash string) error {
	vKey := verifyKey(email)
	cKey := cooldownKey(email)

	pipe := db.Redis.Pipeline()
	pipe.HSet(ctx, vKey, "hash", codeHash, "attempts", 0)
	pipe.Expire(ctx, vKey, codeTTL)
	pipe.Set(ctx, cKey, "1", cooldownTTL)
	_, err := pipe.Exec(ctx)
	return err
}

// GetCode 读取验证码记录
func (s *Store) GetCode(ctx context.Context, email string) (*Record, error) {
	key := verifyKey(email)
	n, err := db.Redis.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, errNotFound
	}

	m, err := db.Redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if m["hash"] == "" {
		return nil, errNotFound
	}

	attempts, _ := strconv.Atoi(m["attempts"])
	return &Record{Hash: m["hash"], Attempts: attempts}, nil
}

// IncrementAttempts 校验失败时递增尝试次数
func (s *Store) IncrementAttempts(ctx context.Context, email string) error {
	return db.Redis.HIncrBy(ctx, verifyKey(email), "attempts", 1).Err()
}

// DeleteCode 校验成功后删除验证码
func (s *Store) DeleteCode(ctx context.Context, email string) error {
	return db.Redis.Del(ctx, verifyKey(email)).Err()
}

// IsNotFound 判断是否为验证码不存在
func IsNotFound(err error) bool {
	return errors.Is(err, errNotFound) || errors.Is(err, redis.Nil)
}
