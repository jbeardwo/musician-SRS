package main

import (
	"container/heap"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jbeardwo/musician-SRS/internal/database"
	"github.com/jbeardwo/musician-SRS/internal/study"
)

func (cfg *apiConfig) nextCardHandler(w http.ResponseWriter, r *http.Request) {
	deckIDstr := r.URL.Query().Get("deck_id")
	deckID, err := uuid.Parse(deckIDstr)
	if err != nil || deckID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Invalid deck ID")
		return
	}

	dbDeck, err := cfg.db.GetDeckById(r.Context(), deckID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Deck not found")
		return
	}

	deck := dbDeckToStudyDeck(dbDeck)

	dbCards, err := cfg.db.GetCardsByDeck(r.Context(), deckID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting cards")
		return
	}
	cardsIn := []study.Card{}
	for _, dbCard := range dbCards {
		cardsIn = append(cardsIn, dbCardToStudyCard(dbCard))
	}
	var cards []study.Card
	var newCards []study.Card

	for _, c := range cardsIn {
		if c.LastReviewedAt.Valid {
			cards = append(cards, c)
		} else {
			newCards = append(newCards, c)
		}
	}

	deck.Cards = cards
	heap.Init(&deck.Cards)
	deck.NewCards = newCards

	nextCard, err := deck.GetNextCard()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting card")
		return
	}

	respondWithJSON(w, http.StatusOK, nextCard)
}

func (cfg *apiConfig) evaluateHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		DeckID uuid.UUID    `json:"deck_id"`
		CardID uuid.UUID    `json:"card_id"`
		Rating study.Rating `json:"rating"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	dbDeck, err := cfg.db.GetDeckById(r.Context(), params.DeckID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid deck ID")
		return
	}
	deck := dbDeckToStudyDeck(dbDeck)

	dbCard, err := cfg.db.GetCardById(r.Context(), params.CardID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid card ID")
		return
	}
	card := dbCardToStudyCard(dbCard)

	if card.DeckID != deck.ID {
		respondWithError(w, http.StatusBadRequest, "Card does not belong to deck")
		return
	}
	rating := params.Rating
	if params.Rating < study.Again || params.Rating > study.Easy {
		respondWithError(w, http.StatusBadRequest, "Invalid rating")
		return
	}
	card.EvaluateCard(rating)
	card.EvaluateMastery(rating)

	deck.TotalReviews++
	card.LastReviewedNum = deck.TotalReviews
	card.Target = deck.TotalReviews + card.Interval

	cardParams := database.UpdateCardParams{
		FrontContent:            card.FrontContent,
		BackContent:             card.BackContent,
		Interval:                card.Interval,
		Target:                  card.Target,
		EaseFactor:              card.EaseFactor,
		RepetitionsCount:        card.RepetitionsCount,
		LastReviewedAt:          card.LastReviewedAt,
		LastReviewedNum:         card.LastReviewedNum,
		Tempo:                   card.Tempo,
		MasteredTempo:           card.MasteredTempo,
		MasteryStage:            string(card.MasteryStage),
		MasteryStageTimeStarted: card.MasteryStageTimeStarted,
		ChallengeAgainCount:     card.ChallengeAgainCount,
		ID:                      card.ID,
	}

	_, err = cfg.db.UpdateCard(r.Context(), cardParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating card")
		return
	}

	deckParams := database.UpdateDeckParams{
		Title:        deck.Title,
		Description:  deck.Description,
		TotalReviews: deck.TotalReviews,
		ID:           deck.ID,
	}

	_, err = cfg.db.UpdateDeck(r.Context(), deckParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating deck")
		return
	}
	respondWithJSON(w, http.StatusOK, card)
}
