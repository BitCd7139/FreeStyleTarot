package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID  `json:"id"`
	Email          string     `json:"email"`
	Nickname       string     `json:"nickname"`
	PasswordHash   *string    `json:"-"`
	Tier           string     `json:"tier"`
	VipExpiresAt   *time.Time `json:"vip_expires_at,omitempty"`
	Balance        float64    `json:"-"`
	LastPredictAt  *time.Time `json:"last_predict_at,omitempty"`
	ExperienceMode *string    `json:"experience_mode,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// PublicUser 返回给前端的用户视图（邮箱已脱敏）
type PublicUser struct {
	ID             uuid.UUID  `json:"id"`
	Email          string     `json:"email"`
	Nickname       string     `json:"nickname"`
	Tier           string     `json:"tier"`
	VipExpiresAt   *time.Time `json:"vip_expires_at,omitempty"`
	Balance        string     `json:"balance"`
	ExperienceMode *string    `json:"experience_mode,omitempty"`
}

func (u *User) ToPublic(maskedEmail string) PublicUser {
	return PublicUser{
		ID:             u.ID,
		Email:          maskedEmail,
		Nickname:       u.Nickname,
		Tier:           u.Tier,
		VipExpiresAt:   u.VipExpiresAt,
		Balance:        FormatBalance(u.Balance),
		ExperienceMode: u.ExperienceMode,
	}
}

func FormatBalance(amount float64) string {
	return fmt.Sprintf("%.2f CNY", amount)
}
