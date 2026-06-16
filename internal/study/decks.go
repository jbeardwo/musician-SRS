package study

import (
	"container/heap"
	"fmt"
	"log"
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

func (d *Deck) ReviewDeck(n int) {
	for i := 0; i < n; {

		if len(d.Cards) == 0 {
			fmt.Println("No Cards in Deck!")
			break
		}

		curCard := heap.Pop(&d.Cards).(Card)
		fmt.Println(curCard.FrontContent)
		fmt.Println(curCard.BackContent)

		fmt.Println("Input: 0. Again, 1. Hard, 2. Good, 3. Easy")

		var evaluation int
		_, err := fmt.Scan(&evaluation)
		if err != nil {
			fmt.Println("Invalid input!")
			heap.Push(&d.Cards, curCard)
			n++
			continue
		}

		log.Println(d.TotalReviews, curCard.Target, curCard.Interval)
		curCard.EvaluateCard(evaluation)
		d.TotalReviews += 1
		curCard.LastReviewedNum = d.TotalReviews
		curCard.Target = d.TotalReviews + curCard.Interval
		log.Println(d.TotalReviews, curCard.Target, curCard.Interval)

		heap.Push(&d.Cards, curCard)
	}
}
