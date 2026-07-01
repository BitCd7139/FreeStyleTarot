package history

import (
	"context"
	"encoding/json"
	"fmt"

	"FreeStyleTarot/db"
	historyModel "FreeStyleTarot/model/predict_history"
	"FreeStyleTarot/model/request"

	"github.com/google/uuid"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Save(ctx context.Context, userID uuid.UUID, limit int, req request.Predict, answer string) error {
	cardSizeJSON, err := json.Marshal(req.CardSize)
	if err != nil {
		return fmt.Errorf("marshal card_size: %w", err)
	}
	cardsJSON, err := json.Marshal(req.Cards)
	if err != nil {
		return fmt.Errorf("marshal cards: %w", err)
	}

	tx, err := db.Pool.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO predict_history (user_id, question, model, card_size, cards, answer)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, req.Question, req.Model, cardSizeJSON, cardsJSON, answer,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM predict_history
		WHERE user_id = $1
		  AND id NOT IN (
		    SELECT id FROM predict_history
		    WHERE user_id = $1
		    ORDER BY created_at DESC
		    LIMIT $2
		  )`, userID, limit)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) ListRecent(ctx context.Context, userID uuid.UUID, limit int) ([]historyModel.Record, error) {
	rows, err := db.Pool.QueryContext(ctx, `
		SELECT id, user_id, question, model, card_size, cards, answer, created_at
		FROM predict_history
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []historyModel.Record
	for rows.Next() {
		var rec historyModel.Record
		var cardSizeJSON, cardsJSON []byte
		if err := rows.Scan(
			&rec.ID, &rec.UserID, &rec.Question, &rec.Model,
			&cardSizeJSON, &cardsJSON, &rec.Answer, &rec.CreatedAt,
		); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(cardSizeJSON, &rec.CardSize); err != nil {
			return nil, fmt.Errorf("unmarshal card_size: %w", err)
		}
		if err := json.Unmarshal(cardsJSON, &rec.Cards); err != nil {
			return nil, fmt.Errorf("unmarshal cards: %w", err)
		}
		records = append(records, rec)
	}
	return records, rows.Err()
}
