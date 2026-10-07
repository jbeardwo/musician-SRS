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
	Target           int32        `json:"target"`
	EaseFactor       float64      `json:"ease_factor"`
	RepetitionsCount int32        `json:"repetitions_count"`
	LastReviewedAt   sql.NullTime `json:"last_reviewed_at"`
	LastReviewedNum  int32        `json:"last_reviewed_num"`
	CreatedAt        time.Time    `json:"created_at"`
	DeckID           uuid.UUID    `json:"deck_id"`
	Tempo            int32        `json:"tempo"`
}

type Rating int

const (
	Again Rating = iota
	Hard
	Good
	Easy
)

func (c *Card) EvaluateCard(r Rating) {
	switch r {
	case Again:
		c.Interval = 1
		c.EaseFactor = math.Max(c.EaseFactor-.20, 1.3)
	case Hard:
		c.Interval = int32(math.Round(float64(c.Interval) * 1.2))
		c.EaseFactor = math.Max(c.EaseFactor-.15, 1.3)
		c.RepetitionsCount += 1
	case Good:
		c.Interval = int32(math.Round(float64(c.Interval) * c.EaseFactor))
		c.RepetitionsCount += 1
	case Easy:
		c.Interval = int32(math.Round(float64(c.Interval) * c.EaseFactor * 1.3))
		c.EaseFactor = c.EaseFactor + .15
		c.RepetitionsCount += 1
	}

	c.LastReviewedAt = sql.NullTime{Time: time.Now(), Valid: true}
}
