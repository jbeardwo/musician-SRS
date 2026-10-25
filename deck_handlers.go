package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jbeardwo/musician-SRS/internal/database"
	"github.com/jbeardwo/musician-SRS/internal/study"
)

func (cfg *apiConfig) getDecksHandler(w http.ResponseWriter, r *http.Request) {
	userIDstr := r.URL.Query().Get("user_id")

	userID, err := uuid.Parse(userIDstr)
	if err != nil || userID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	dbDecks, err := cfg.db.GetDecksByUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting decks")
		return
	}

	decks := []study.Deck{}
	for _, dbDeck := range dbDecks {
		decks = append(decks, dbDeckToStudyDeck(dbDeck))
	}
	respondWithJSON(w, http.StatusOK, decks)
}

func (cfg *apiConfig) getDeckHandler(w http.ResponseWriter, r *http.Request) {

	deckIDstr := r.PathValue("deckID")
	deckID, err := uuid.Parse(deckIDstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid deck ID")
		return
	}

	dbDeck, err := cfg.db.GetDeckById(r.Context(), deckID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Deck not found")
		return
	}

	deck := dbDeckToStudyDeck(dbDeck)

	respondWithJSON(w, http.StatusOK, deck)
}

func (cfg *apiConfig) deleteDeckHandler(w http.ResponseWriter, r *http.Request) {
	deckIDstr := r.PathValue("deckID")
	deckID, err := uuid.Parse(deckIDstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid deck ID")
		return
	}

	_, err = cfg.db.GetDeckById(r.Context(), deckID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Deck not found")
		return
	}

	err = cfg.db.DeleteDeck(r.Context(), deckID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) newDeckHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		UserID      uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	deckParams := database.CreateDeckParams{
		Title:       params.Title,
		Description: params.Description,
		UserID:      params.UserID,
	}

	dbDeck, err := cfg.db.CreateDeck(r.Context(), deckParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating deck")
		return
	}

	deck := dbDeckToStudyDeck(dbDeck)

	respondWithJSON(w, http.StatusCreated, deck)
}

func (cfg *apiConfig) updateDeckHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		TotalReviews int32     `json:"total_reviews"`
		ID           uuid.UUID `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	deckParams := database.UpdateDeckParams{
		Title:        params.Title,
		Description:  params.Description,
		TotalReviews: params.TotalReviews,
		ID:           params.ID,
	}

	dbDeck, err := cfg.db.UpdateDeck(r.Context(), deckParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating deck")
		return
	}

	deck := dbDeckToStudyDeck(dbDeck)

	respondWithJSON(w, http.StatusOK, deck)
}
