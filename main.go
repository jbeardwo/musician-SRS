package main

import (
	"container/heap"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jbeardwo/musician-SRS/internal/database"
	"github.com/jbeardwo/musician-SRS/internal/study"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db  *database.Queries
	ctx context.Context
}

func main() {
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
		ID:           dbDeck.ID,
		Title:        dbDeck.Title,
		Description:  dbDeck.Description,
		CreatedAt:    dbDeck.CreatedAt,
		UserID:       dbDeck.UserID,
		TotalReviews: dbDeck.TotalReviews,
	}

	for _, note := range study.CommonNotes {
		cardParams := database.CreateCardParams{
			FrontContent: note + " Major Scale",
			BackContent:  "DO IT",
			Target:       1,
			DeckID:       deck.ID,
			Tempo:        50,
		}
		dbCard, err := dbQueries.CreateCard(ctx, cardParams)
		if err != nil {
			log.Fatal(err)
		}

		card := study.Card{
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

		heap.Push(&deck.Cards, card)
	}

	apiCfg := apiConfig{
		db:  dbQueries,
		ctx: ctx,
	}
	// deck.ReviewDeck(4)
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("GET /api/healthz", readyHandler)
	serveMux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	serveMux.HandleFunc("GET /api/decks", apiCfg.getDecksHandler)
	serveMux.HandleFunc("GET /api/decks/{deckID}", apiCfg.getDeckHandler)
	serveMux.HandleFunc("GET /api/cards", apiCfg.getCardsHandler)
	serveMux.HandleFunc("GET /api/cards/{cardID}", apiCfg.getCardHandler)
	serveMux.HandleFunc("POST /api/decks", apiCfg.newDeckHandler)
	serveMux.HandleFunc("POST /api/cards", apiCfg.newCardHandler)
	serveMux.HandleFunc("DELETE /api/decks/{deckID}", apiCfg.deleteDeckHandler)
	serveMux.HandleFunc("DELETE /api/cards/{cardID}", apiCfg.deleteCardHandler)
	serveMux.HandleFunc("PUT /api/decks", apiCfg.updateDeckHandler)
	serveMux.HandleFunc("PUT /api/cards", apiCfg.updateCardHandler)
	serveMux.HandleFunc("GET /{$}", homeHandler)
	serveMux.HandleFunc("GET /js/login.js", loginJSHandler)
	serveMux.HandleFunc("GET /decks", decksHandler)
	serveMux.HandleFunc("GET /js/decks.js", decksJSHandler)
	serveMux.HandleFunc("GET /edit-deck", editDeckHandler)
	serveMux.HandleFunc("GET /js/edit-deck.js", editDeckJSHandler)
	serveMux.HandleFunc("GET /edit-card", editCardHandler)
	serveMux.HandleFunc("GET /js/edit-card.js", editCardJSHandler)
	serveMux.HandleFunc("GET /study", studyHandler)
	serveMux.HandleFunc("GET /js/study.js", studyJSHandler)
	serveMux.HandleFunc("GET /api/study/next", apiCfg.nextCardHandler)
	serveMux.HandleFunc("POST /api/study/evaluate", apiCfg.evaluateHandler)

	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func loginJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/login.js")
}

func decksJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/decks.js")
}

func editDeckJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/edit-deck.js")
}

func editCardJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/edit-card.js")
}

func studyJSHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/js/study.js")
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding body")
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	type loginResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	response := loginResponse{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/index.html")
}

func decksHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/decks.html")
}

func editDeckHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/edit-deck.html")
}

func editCardHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/edit-card.html")
}

func studyHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/html/study.html")
}

func (cfg *apiConfig) getDecksHandler(w http.ResponseWriter, r *http.Request) {
	userIDstr := r.URL.Query().Get("user_id")

	var dbDecks []database.Deck
	userID, err := uuid.Parse(userIDstr)
	if err != nil || userID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	dbDecks, err = cfg.db.GetDecksByUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, 500, "Error getting decks")
		return
	}

	decks := []study.Deck{}
	for _, dbDeck := range dbDecks {
		decks = append(decks, study.Deck{
			ID:           dbDeck.ID,
			Title:        dbDeck.Title,
			Description:  dbDeck.Description,
			CreatedAt:    dbDeck.CreatedAt,
			UserID:       dbDeck.UserID,
			TotalReviews: dbDeck.TotalReviews,
		})
	}
	respondWithJSON(w, http.StatusOK, decks)
}

func (cfg *apiConfig) deleteCardHandler(w http.ResponseWriter, r *http.Request) {
	cardIDstr := r.PathValue("cardID")
	cardID, err := uuid.Parse(cardIDstr)
	if err != nil {
		respondWithError(w, 400, "Invalid card ID")
		return
	}

	_, err = cfg.db.GetCardById(r.Context(), cardID)
	if err != nil {
		respondWithError(w, 404, "Card not found")
		return
	}

	err = cfg.db.DeleteCard(r.Context(), cardID)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	w.WriteHeader(204)
}

func (cfg *apiConfig) nextCardHandler(w http.ResponseWriter, r *http.Request) {
	deckIDstr := r.URL.Query().Get("deck_id")
	deckID, err := uuid.Parse(deckIDstr)
	if err != nil || deckID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Invalid deck ID")
		return
	}

	dbDeck, err := cfg.db.GetDeckById(r.Context(), deckID)
	if err != nil {
		respondWithError(w, 404, "Deck not found")
		return
	}

	deck := study.Deck{
		ID:           dbDeck.ID,
		Title:        dbDeck.Title,
		Description:  dbDeck.Description,
		CreatedAt:    dbDeck.CreatedAt,
		UserID:       dbDeck.UserID,
		TotalReviews: dbDeck.TotalReviews,
	}

	var dbCards []database.Card
	dbCards, err = cfg.db.GetCardsByDeck(r.Context(), deckID)
	if err != nil {
		respondWithError(w, 500, "Error getting cards")
		return
	}
	cardsIn := []study.Card{}
	for _, dbCard := range dbCards {
		cardsIn = append(cardsIn, study.Card{
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
		})
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
		respondWithError(w, 500, "error getting card")
		return
	}

	respondWithJSON(w, http.StatusOK, nextCard)
}

func (cfg *apiConfig) getCardsHandler(w http.ResponseWriter, r *http.Request) {
	deckIDstr := r.URL.Query().Get("deck_id")

	var dbCards []database.Card
	deckID, err := uuid.Parse(deckIDstr)
	if err != nil || deckID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	dbCards, err = cfg.db.GetCardsByDeck(r.Context(), deckID)
	if err != nil {
		respondWithError(w, 500, "Error getting decks")
		return
	}

	cards := []study.Card{}
	for _, dbCard := range dbCards {
		cards = append(cards, study.Card{
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
		})
	}

	respondWithJSON(w, http.StatusOK, cards)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errReturnVals struct {
		Error string `json:"error"`
	}
	respBody := errReturnVals{
		Error: msg,
	}
	dat, _ := json.Marshal(respBody)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(dat)
}

func (cfg *apiConfig) getDeckHandler(w http.ResponseWriter, r *http.Request) {

	deckIDstr := r.PathValue("deckID")
	deckID, err := uuid.Parse(deckIDstr)
	if err != nil {
		respondWithError(w, 400, "Invalid deck ID")
		return
	}

	dbDeck, err := cfg.db.GetDeckById(r.Context(), deckID)
	if err != nil {
		respondWithError(w, 404, "Deck not found")
		return
	}

	deck := study.Deck{
		ID:           dbDeck.ID,
		Title:        dbDeck.Title,
		Description:  dbDeck.Description,
		CreatedAt:    dbDeck.CreatedAt,
		UserID:       dbDeck.UserID,
		TotalReviews: dbDeck.TotalReviews,
	}

	respondWithJSON(w, http.StatusOK, deck)
}

func (cfg *apiConfig) getCardHandler(w http.ResponseWriter, r *http.Request) {

	cardIDstr := r.PathValue("cardID")
	cardID, err := uuid.Parse(cardIDstr)
	if err != nil {
		respondWithError(w, 400, "Invalid card ID")
		return
	}

	dbCard, err := cfg.db.GetCardById(r.Context(), cardID)
	if err != nil {
		respondWithError(w, 404, "Card not found")
		return
	}

	card := study.Card{
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

	respondWithJSON(w, http.StatusOK, card)
}

func (cfg *apiConfig) deleteDeckHandler(w http.ResponseWriter, r *http.Request) {
	deckIDstr := r.PathValue("deckID")
	deckID, err := uuid.Parse(deckIDstr)
	if err != nil {
		respondWithError(w, 400, "Invalid deck ID")
		return
	}

	_, err = cfg.db.GetDeckById(r.Context(), deckID)
	if err != nil {
		respondWithError(w, 404, "Deck not found")
		return
	}

	err = cfg.db.DeleteDeck(r.Context(), deckID)
	if err != nil {
		respondWithError(w, 500, "Something went wong")
		return
	}

	w.WriteHeader(204)
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

	dbDeck, err := cfg.db.CreateDeck(cfg.ctx, deckParams)
	if err != nil {
		log.Fatal(err)
	}

	deck := study.Deck{
		ID:           dbDeck.ID,
		Title:        dbDeck.Title,
		Description:  dbDeck.Description,
		CreatedAt:    dbDeck.CreatedAt,
		UserID:       dbDeck.UserID,
		TotalReviews: dbDeck.TotalReviews,
	}

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

	dbDeck, err := cfg.db.UpdateDeck(cfg.ctx, deckParams)
	if err != nil {
		log.Fatal(err)
	}

	deck := study.Deck{
		ID:           dbDeck.ID,
		Title:        dbDeck.Title,
		Description:  dbDeck.Description,
		CreatedAt:    dbDeck.CreatedAt,
		UserID:       dbDeck.UserID,
		TotalReviews: dbDeck.TotalReviews,
	}

	respondWithJSON(w, http.StatusOK, deck)
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

	dbDeck, err := cfg.db.GetDeckById(cfg.ctx, params.DeckID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Deck ID")
		return
	}
	deck := study.Deck{
		ID:           dbDeck.ID,
		Title:        dbDeck.Title,
		Description:  dbDeck.Description,
		CreatedAt:    dbDeck.CreatedAt,
		UserID:       dbDeck.UserID,
		TotalReviews: dbDeck.TotalReviews,
	}

	dbCard, err := cfg.db.GetCardById(cfg.ctx, params.CardID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Card ID")
		return
	}
	card := study.Card{
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

	dbCard, err := cfg.db.UpdateCard(cfg.ctx, cardParams)
	if err != nil {
		log.Fatal(err)
	}

	card := study.Card{
		ID:                      dbCard.ID,
		FrontContent:            dbCard.FrontContent,
		BackContent:             dbCard.BackContent,
		Interval:                dbCard.Interval,
		Target:                  dbCard.Target,
		EaseFactor:              dbCard.EaseFactor,
		RepetitionsCount:        dbCard.RepetitionsCount,
		LastReviewedAt:          dbCard.LastReviewedAt,
		LastReviewedNum:         dbCard.LastReviewedNum,
		DeckID:                  dbCard.DeckID,
		CreatedAt:               dbCard.CreatedAt,
		Tempo:                   dbCard.Tempo,
		MasteredTempo:           dbCard.MasteredTempo,
		MasteryStage:            study.MasteryStage(dbCard.MasteryStage),
		MasteryStageTimeStarted: dbCard.MasteryStageTimeStarted,
		ChallengeAgainCount:     dbCard.ChallengeAgainCount,
	}

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

	dbCard, err := cfg.db.CreateCard(cfg.ctx, cardParams)
	if err != nil {
		log.Fatal(err)
	}

	card := study.Card{
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

	respondWithJSON(w, http.StatusCreated, card)
}
