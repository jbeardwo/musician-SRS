package main

import (
	"github.com/jbeardwo/musician-SRS/internal/database"
	"github.com/jbeardwo/musician-SRS/internal/study"
)

func dbDeckToStudyDeck(dbDeck database.Deck) study.Deck {
	return study.Deck{
		ID:           dbDeck.ID,
		Title:        dbDeck.Title,
		Description:  dbDeck.Description,
		CreatedAt:    dbDeck.CreatedAt,
		UserID:       dbDeck.UserID,
		TotalReviews: dbDeck.TotalReviews,
	}
}

func dbCardToStudyCard(dbCard database.Card) study.Card {
	return study.Card{
		ID:                      dbCard.ID,
		FrontContent:            dbCard.FrontContent,
		BackContent:             dbCard.BackContent,
		Interval:                dbCard.Interval,
		Target:                  dbCard.Target,
		EaseFactor:              dbCard.EaseFactor,
		RepetitionsCount:        dbCard.RepetitionsCount,
		LastReviewedAt:          dbCard.LastReviewedAt,
		LastReviewedNum:         dbCard.LastReviewedNum,
		CreatedAt:               dbCard.CreatedAt,
		DeckID:                  dbCard.DeckID,
		Tempo:                   dbCard.Tempo,
		MasteredTempo:           dbCard.MasteredTempo,
		MasteryStage:            study.MasteryStage(dbCard.MasteryStage),
		MasteryStageTimeStarted: dbCard.MasteryStageTimeStarted,
		ChallengeAgainCount:     dbCard.ChallengeAgainCount,
	}
}

func dbUserToStudyUser(dbUser database.User) study.User {
	return study.User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
}
