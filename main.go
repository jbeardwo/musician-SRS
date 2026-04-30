package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/jbeardwo/musician-SRS/internal/database"
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

	dbUser, err := dbQueries.CreateUser(ctx, "wilford@brimley.com")
	if err != nil {
		log.Fatal(err)
	}

	user := User{
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

	deck := Deck{
		ID:          dbDeck.ID,
		Title:       dbDeck.Title,
		Description: dbDeck.Description,
		CreatedAt:   dbDeck.CreatedAt,
		UserID:      dbDeck.UserID,
	}

	log.Println(user.ID)
	log.Println(deck.ID)
	log.Println(deck.UserID)
	log.Println(deck.Title)
	log.Println(deck.Description)
}
