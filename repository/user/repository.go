package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"FreeStyleTarot/db"
	userModel "FreeStyleTarot/model/user"

	"github.com/google/uuid"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func scanUser(row interface {
	Scan(dest ...any) error
}) (*userModel.User, error) {
	var u userModel.User
	err := row.Scan(
		&u.ID, &u.Email, &u.Nickname, &u.PasswordHash, &u.Tier, &u.VipExpiresAt,
		&u.Balance, &u.LastPredictAt, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

const userColumns = `id, email, nickname, password_hash, tier, vip_expires_at,
	balance, last_predict_at, created_at, updated_at`

// GetByEmail 按邮箱查询用户
func (r *Repository) GetByEmail(ctx context.Context, email string) (*userModel.User, error) {
	row := db.Pool.QueryRowContext(ctx,
		`SELECT `+userColumns+` FROM users WHERE email = $1`, email)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

// GetByID 按 ID 查询用户
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*userModel.User, error) {
	row := db.Pool.QueryRowContext(ctx,
		`SELECT `+userColumns+` FROM users WHERE id = $1`, id)
	u, err := scanUser(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return u, err
}

// CreateWithPassword 邮箱密码注册
func (r *Repository) CreateWithPassword(ctx context.Context, email, nickname, passwordHash string) (*userModel.User, error) {
	row := db.Pool.QueryRowContext(ctx, `
		INSERT INTO users (email, nickname, password_hash)
		VALUES ($1, $2, $3)
		RETURNING `+userColumns,
		email, nickname, passwordHash,
	)
	return scanUser(row)
}

// EnsureByEmail 确保用户存在（验证码登录成功后创建）
func (r *Repository) EnsureByEmail(ctx context.Context, email string) (*userModel.User, error) {
	u, err := r.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return u, nil
	}
	return r.Create(ctx, email, "")
}

// Create 创建新用户
func (r *Repository) Create(ctx context.Context, email, nickname string) (*userModel.User, error) {
	row := db.Pool.QueryRowContext(ctx, `
		INSERT INTO users (email, nickname)
		VALUES ($1, $2)
		RETURNING `+userColumns,
		email, nickname,
	)
	return scanUser(row)
}

// UpdateNickname 更新昵称
func (r *Repository) UpdateNickname(ctx context.Context, id uuid.UUID, nickname string) error {
	_, err := db.Pool.ExecContext(ctx, `
		UPDATE users SET nickname = $2, updated_at = NOW() WHERE id = $1`,
		id, nickname)
	return err
}

// UpdateLastPredictAt 更新上次占卜时间
func (r *Repository) UpdateLastPredictAt(ctx context.Context, id uuid.UUID, at time.Time) error {
	_, err := db.Pool.ExecContext(ctx, `
		UPDATE users SET last_predict_at = $2, updated_at = NOW() WHERE id = $1`,
		id, at)
	return err
}

// UpdatePassword 设置密码
func (r *Repository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	_, err := db.Pool.ExecContext(ctx, `
		UPDATE users SET password_hash = $2, updated_at = NOW() WHERE id = $1`,
		id, passwordHash)
	return err
}

// SetNicknameIfEmpty 注册时若提供了昵称则写入
func (r *Repository) SetNicknameIfEmpty(ctx context.Context, email, nickname string) error {
	_, err := db.Pool.ExecContext(ctx, `
		UPDATE users SET nickname = $2, updated_at = NOW()
		WHERE email = $1 AND nickname = ''`, email, nickname)
	return err
}
