package study

import (
	"database/sql"
	"math"
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

func (c *Card) EvaluateCard(eval int) {
	if eval == 0 {
		c.Interval = 1
		c.EaseFactor = math.Max(c.EaseFactor*.80, 1.3)
		c.RepetitionsCount = 0
	} else if eval <= 3 {
		c.Interval = int32(math.Round(float64(c.Interval) * c.EaseFactor))
		c.EaseFactor = math.Max(c.EaseFactor*(0.7+(0.15*float64(eval))), 1.3)
		c.RepetitionsCount++
	}

	c.LastReviewedAt = sql.NullTime{Time: time.Now(), Valid: true}
}
