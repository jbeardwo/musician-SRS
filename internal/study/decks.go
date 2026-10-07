package study

import (
	"time"

	"github.com/google/uuid"
)

type Deck struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UserID       uuid.UUID `json:"user_id"`
	TotalReviews int32     `json:"total_reviews"`
	Cards        CardHeap
}
