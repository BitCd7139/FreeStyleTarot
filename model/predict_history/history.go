package predict_history

import (
	"time"

	"FreeStyleTarot/model/request"

	"github.com/google/uuid"
)

type Record struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Question  string
	Model     string
	CardSize  request.CardSize
	Cards     []request.CardInfo
	Answer    string
	CreatedAt time.Time
}

type Item struct {
	ID        uuid.UUID          `json:"id"`
	Question  string             `json:"question"`
	Model     string             `json:"model"`
	CardSize  request.CardSize   `json:"cardSize"`
	Cards     []request.CardInfo `json:"cards"`
	Answer    string             `json:"answer"`
	CreatedAt time.Time          `json:"created_at"`
}

func (r *Record) ToItem() Item {
	return Item{
		ID:        r.ID,
		Question:  r.Question,
		Model:     r.Model,
		CardSize:  r.CardSize,
		Cards:     r.Cards,
		Answer:    r.Answer,
		CreatedAt: r.CreatedAt,
	}
}

type ListResult struct {
	Limit int    `json:"limit"`
	Items []Item `json:"items"`
}
