package study

import (
	"container/heap"
	"errors"
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
	NewCards     []Card
}

func (d *Deck) GetNextCard() (Card, error) {
	card, err := d.selectNextCard()
	if err != nil {
		return Card{}, err
	}
	// if this card is the one most recently reviewed, get another
	if card.LastReviewedAt.Valid && card.LastReviewedNum == d.TotalReviews {
		otherCard, err := d.selectNextCard()
		//if there are no other cards then give the repeat
		if err != nil {
			return card, nil
		}
		heap.Push(&d.Cards, card)
		return otherCard, nil
	}
	return card, nil
}

func (d *Deck) selectNextCard() (Card, error) {
	if len(d.Cards) == 0 && len(d.NewCards) == 0 {
		return Card{}, errors.New("no card available")
	}

	if len(d.Cards) == 0 {
		card := d.NewCards[0]
		d.NewCards = d.NewCards[1:]
		return card, nil
	}

	root := d.Cards[0]

	if root.Target <= d.TotalReviews {
		return heap.Pop(&d.Cards).(Card), nil
	}
	if len(d.NewCards) > 0 {
		card := d.NewCards[0]
		d.NewCards = d.NewCards[1:]
		return card, nil
	}
	return heap.Pop(&d.Cards).(Card), nil
}
