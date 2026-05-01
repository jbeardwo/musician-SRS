package main

import (
	"container/heap"
	"context"
	"database/sql"
	"log"

	"github.com/jbeardwo/musician-SRS/internal/database"
	"github.com/jbeardwo/musician-SRS/internal/study"
	_ "github.com/lib/pq"
)

func main() {
	const fileRoot = "."
	const port = "8080"
	const dbURL = "postgres://postgres@localhost:5432/musiciansrs?sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	ctx := context.Background()

	err = dbQueries.DeleteUsers(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = dbQueries.DeleteDecks(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = dbQueries.DeleteCards(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dbUser, err := dbQueries.CreateUser(ctx, "wilford@brimley.com")
	if err != nil {
		log.Fatal(err)
	}

	user := study.User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	deckParams := database.CreateDeckParams{
		Title:       "Deck",
		Description: "Testing The Program",
		UserID:      user.ID,
	}

	dbDeck, err := dbQueries.CreateDeck(ctx, deckParams)
	if err != nil {
		log.Fatal(err)
	}

	deck := study.Deck{
		ID:          dbDeck.ID,
		Title:       dbDeck.Title,
		Description: dbDeck.Description,
		CreatedAt:   dbDeck.CreatedAt,
		UserID:      dbDeck.UserID,
	}

	for _, note := range study.CommonNotes {
		cardParams := database.CreateCardParams{
			FrontContent: note + " Major Scale",
			BackContent:  "DO IT",
			DeckID:       deck.ID,
		}
		dbCard, err := dbQueries.CreateCard(ctx, cardParams)
		if err != nil {
			log.Fatal(err)
		}

		card := study.Card{
			ID:               dbCard.ID,
			FrontContent:     dbCard.FrontContent,
			BackContent:      dbCard.BackContent,
			Interval:         dbCard.Interval,
			EaseFactor:       dbCard.EaseFactor,
			RepetitionsCount: dbCard.RepetitionsCount,
			LastReviewedAt:   dbCard.LastReviewedAt,
			CreatedAt:        dbCard.CreatedAt,
			DeckID:           dbCard.DeckID,
		}

		heap.Push(&deck.Cards, card)
	}

	deck.ReviewDeck(4)
}
