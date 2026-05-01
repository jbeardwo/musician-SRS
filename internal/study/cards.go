package study

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID               uuid.UUID    `json:"id"`
	FrontContent     string       `json:"front_content"`
	BackContent      string       `json:"back_content"`
	Interval         int32        `json:"interval"`
	EaseFactor       float64      `json:"ease_factor"`
	RepetitionsCount int32        `json:"repetitions_count"`
	LastReviewedAt   sql.NullTime `json:"last_reviewed_at"`
	CreatedAt        time.Time    `json:"created_at"`
	DeckID           uuid.UUID    `json:"deck_id"`
}
