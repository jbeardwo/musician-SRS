package study

import (
	"database/sql"
	"math"
	"time"

	"github.com/google/uuid"
)

type MasteryStage string

const (
	Challenge MasteryStage = "challenge"
	Playable  MasteryStage = "playable"
	Retained  MasteryStage = "retained"
)

type Card struct {
	ID                      uuid.UUID    `json:"id"`
	FrontContent            string       `json:"front_content"`
	BackContent             string       `json:"back_content"`
	Interval                int32        `json:"interval"`
	Target                  int32        `json:"target"`
	EaseFactor              float64      `json:"ease_factor"`
	RepetitionsCount        int32        `json:"repetitions_count"`
	LastReviewedAt          sql.NullTime `json:"last_reviewed_at"`
	LastReviewedNum         int32        `json:"last_reviewed_num"`
	CreatedAt               time.Time    `json:"created_at"`
	DeckID                  uuid.UUID    `json:"deck_id"`
	Tempo                   int32        `json:"tempo"`
	MasteredTempo           int32        `json:"mastered_tempo"`
	MasteryStage            MasteryStage `json:"mastery_stage"`
	MasteryStageTimeStarted sql.NullTime `json:"mastery_stage_time_started"`
	ChallengeAgainCount     int32        `json:"challenge_again_count"`
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

func (c *Card) EvaluateMastery(r Rating) {

	switch c.MasteryStage {
	case Challenge:
		switch r {
		case Again:
			c.ChallengeAgainCount++
			if c.ChallengeAgainCount == 3 {
				c.Tempo -= 2
				c.ChallengeAgainCount = 0
			}
		case Hard:
		case Good, Easy:
			c.MasteryStage = Playable
			c.ChallengeAgainCount = 0
			c.MasteryStageTimeStarted = sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}
		}
	case Playable:
		if !c.MasteryStageTimeStarted.Valid && r != Again {
			c.MasteryStageTimeStarted = sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}
			return
		}

		switch r {
		case Again:
			c.MasteryStage = Challenge
			c.ChallengeAgainCount = 0
			c.MasteryStageTimeStarted = sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}
		case Hard:
			if c.MasteryStageTimeStarted.Valid && time.Since(c.MasteryStageTimeStarted.Time) >= 8*time.Hour {
				c.MasteryStageTimeStarted = sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				}
			}
		case Good, Easy:
			if c.MasteryStageTimeStarted.Valid && time.Since(c.MasteryStageTimeStarted.Time) >= 8*time.Hour {
				c.MasteryStage = Retained
				c.MasteryStageTimeStarted = sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				}
			}
		}
	case Retained:

		switch r {
		case Again:
			c.MasteryStage = Playable
			c.MasteryStageTimeStarted.Time = time.Now()
		case Hard:
			if c.MasteryStageTimeStarted.Valid && time.Since(c.MasteryStageTimeStarted.Time) >= 20*time.Hour {
				c.MasteryStageTimeStarted.Time = time.Now()
			}
		case Good, Easy:
			if c.MasteryStageTimeStarted.Valid && time.Since(c.MasteryStageTimeStarted.Time) >= 20*time.Hour {
				c.MasteredTempo = c.Tempo
				c.Tempo += 5
				c.MasteryStage = Challenge
				c.ChallengeAgainCount = 0
				c.MasteryStageTimeStarted.Time = time.Now()
			}
		}

	}
}
