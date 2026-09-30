package study

import (
	"time"

	"github.com/google/uuid"
)

type Deck struct {
	ID               uuid.UUID `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UserID           uuid.UUID `json:"user_id"`
	TotalReviews     int32     `json:"total_reviews"`
	TempoIntervalUp  int32     `json:"tempo_interval_up"`
	TempoIntervalDn  int32     `json:"tempo_interval_dn"`
	PerfectThreshold int32     `json:"perfect_threshold"`
	BadThreshold     int32     `json:"bad_threshold"`
	Cards            CardHeap
}
