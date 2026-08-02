package study

import (
	"container/heap"
	"fmt"
	"log"
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

func (d *Deck) ReviewDeck(n int) {
	for range n {
		if len(d.Cards) < 4 {
			fmt.Println("4 cards minimum required to study")
			break
		}

		curCard := heap.Pop(&d.Cards).(Card)

		// enforce at least 3 cards before repetition
		if d.TotalReviews-curCard.LastReviewedNum < 3 && d.TotalReviews >= 3 {
			var heldCards []Card
			heldCards = append(heldCards, curCard)
			for {
				curCard = heap.Pop(&d.Cards).(Card)
				if d.TotalReviews-curCard.LastReviewedNum < 3 {
					heldCards = append(heldCards, curCard)
				} else {
					break
				}
			}
			for _, card := range heldCards {
				heap.Push(&d.Cards, card)
			}
		}

		fmt.Println(curCard.FrontContent)
		fmt.Println(curCard.Tempo)
		fmt.Println(curCard.BackContent)

		fmt.Println("Input: 0. Again, 1. Hard, 2. Good, 3. Easy")

		var evaluation int

		for {
			_, err := fmt.Scan(&evaluation)
			if err != nil {
				fmt.Println("Not an integer!")
			}
			if evaluation >= 0 && evaluation <= 3 {
				break
			}
			fmt.Println("Invalid input: evaluation must be an integer between 0 and 3")
		}

		log.Println(d.TotalReviews, curCard.Target, curCard.Interval, curCard.Tempo)
		curCard.EvaluateCard(evaluation)
		d.TotalReviews += 1
		curCard.LastReviewedNum = d.TotalReviews
		curCard.Target = d.TotalReviews + curCard.Interval
		if curCard.BadStreak == d.BadThreshold {
			curCard.Tempo -= d.TempoIntervalDn
			curCard.BadStreak = 0
		}
		if curCard.PerfectStreak == d.PerfectThreshold {
			curCard.Tempo += d.TempoIntervalUp
			curCard.PerfectStreak = 0
		}
		log.Println(d.TotalReviews, curCard.Target, curCard.Interval, curCard.Tempo)

		heap.Push(&d.Cards, curCard)
	}
}
