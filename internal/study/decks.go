package study

import (
	"time"

	"github.com/google/uuid"
)

type Deck struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UserID      uuid.UUID `json:"user_id"`
	Cards       CardHeap
}

func (d *Deck) reviewDeck(n int) {
}
