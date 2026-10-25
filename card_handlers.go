package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jbeardwo/musician-SRS/internal/database"
	"github.com/jbeardwo/musician-SRS/internal/study"
)

func (cfg *apiConfig) getCardsHandler(w http.ResponseWriter, r *http.Request) {
	deckIDstr := r.URL.Query().Get("deck_id")

	deckID, err := uuid.Parse(deckIDstr)
	if err != nil || deckID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Invalid deck ID")
		return
	}
	dbCards, err := cfg.db.GetCardsByDeck(r.Context(), deckID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting cards")
		return
	}

	cards := []study.Card{}
	for _, dbCard := range dbCards {
		cards = append(cards, dbCardToStudyCard(dbCard))
	}

	respondWithJSON(w, http.StatusOK, cards)
}

func (cfg *apiConfig) getCardHandler(w http.ResponseWriter, r *http.Request) {

	cardIDstr := r.PathValue("cardID")
	cardID, err := uuid.Parse(cardIDstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid card ID")
		return
	}

	dbCard, err := cfg.db.GetCardById(r.Context(), cardID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Card not found")
		return
	}

	card := dbCardToStudyCard(dbCard)

	respondWithJSON(w, http.StatusOK, card)
}

func (cfg *apiConfig) newCardHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FrontContent string    `json:"front_content"`
		BackContent  string    `json:"back_content"`
		Target       int32     `json:"target"`
		DeckID       uuid.UUID `json:"deck_id"`
		Tempo        int32     `json:"tempo"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	cardParams := database.CreateCardParams{
		FrontContent: params.FrontContent,
		BackContent:  params.BackContent,
		Target:       params.Target,
		DeckID:       params.DeckID,
		Tempo:        params.Tempo,
	}

	dbCard, err := cfg.db.CreateCard(r.Context(), cardParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating card")
		return
	}

	card := dbCardToStudyCard(dbCard)
	respondWithJSON(w, http.StatusCreated, card)
}

func (cfg *apiConfig) deleteCardHandler(w http.ResponseWriter, r *http.Request) {
	cardIDstr := r.PathValue("cardID")
	cardID, err := uuid.Parse(cardIDstr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid card ID")
		return
	}

	_, err = cfg.db.GetCardById(r.Context(), cardID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Card not found")
		return
	}

	err = cfg.db.DeleteCard(r.Context(), cardID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) updateCardHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		FrontContent            string       `json:"front_content"`
		BackContent             string       `json:"back_content"`
		Interval                int32        `json:"interval"`
		Target                  int32        `json:"target"`
		EaseFactor              float64      `json:"ease_factor"`
		RepetitionsCount        int32        `json:"repetitions_count"`
		LastReviewedAt          sql.NullTime `json:"last_reviewed_at"`
		LastReviewedNum         int32        `json:"last_reviewed_num"`
		Tempo                   int32        `json:"tempo"`
		ID                      uuid.UUID    `json:"id"`
		MasteredTempo           int32        `json:"mastered_tempo"`
		MasteryStage            string       `json:"mastery_stage"`
		MasteryStageTimeStarted sql.NullTime `json:"mastery_stage_time_started"`
		ChallengeAgainCount     int32        `json:"challenge_again_count"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	cardParams := database.UpdateCardParams{
		FrontContent:            params.FrontContent,
		BackContent:             params.BackContent,
		Interval:                params.Interval,
		Target:                  params.Target,
		EaseFactor:              params.EaseFactor,
		RepetitionsCount:        params.RepetitionsCount,
		LastReviewedAt:          params.LastReviewedAt,
		LastReviewedNum:         params.LastReviewedNum,
		Tempo:                   params.Tempo,
		ID:                      params.ID,
		MasteredTempo:           params.MasteredTempo,
		MasteryStage:            params.MasteryStage,
		MasteryStageTimeStarted: params.MasteryStageTimeStarted,
		ChallengeAgainCount:     params.ChallengeAgainCount,
	}

	dbCard, err := cfg.db.UpdateCard(r.Context(), cardParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating card")
		return
	}

	card := dbCardToStudyCard(dbCard)

	respondWithJSON(w, http.StatusOK, card)
}
