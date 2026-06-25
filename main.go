package main

import (
	"container/heap"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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
		Title:            "Deck",
		Description:      "Testing The Program",
		UserID:           user.ID,
		TempoIntervalUp:  5,
		TempoIntervalDn:  2,
		PerfectThreshold: 5,
		BadThreshold:     2,
	}

	dbDeck, err := dbQueries.CreateDeck(ctx, deckParams)
	if err != nil {
		log.Fatal(err)
	}

	deck := study.Deck{
		ID:               dbDeck.ID,
		Title:            dbDeck.Title,
		Description:      dbDeck.Description,
		CreatedAt:        dbDeck.CreatedAt,
		UserID:           dbDeck.UserID,
		TotalReviews:     dbDeck.TotalReviews,
		TempoIntervalUp:  dbDeck.TempoIntervalUp,
		TempoIntervalDn:  dbDeck.TempoIntervalDn,
		PerfectThreshold: dbDeck.PerfectThreshold,
		BadThreshold:     dbDeck.BadThreshold,
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
			ID:               dbCard.ID,
			FrontContent:     dbCard.FrontContent,
			BackContent:      dbCard.BackContent,
			Interval:         dbCard.Interval,
			Target:           dbCard.Target,
			EaseFactor:       dbCard.EaseFactor,
			RepetitionsCount: dbCard.RepetitionsCount,
			LastReviewedAt:   dbCard.LastReviewedAt,
			CreatedAt:        dbCard.CreatedAt,
			DeckID:           dbCard.DeckID,
			Tempo:            dbCard.Tempo,
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
	serveMux.HandleFunc("GET /api/cards", apiCfg.getCardsHandler)
	serveMux.HandleFunc("POST /api/decks", apiCfg.newDeckHandler)

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
			ID:          dbDeck.ID,
			Title:       dbDeck.Title,
			Description: dbDeck.Description,
			CreatedAt:   dbDeck.CreatedAt,
			UserID:      dbDeck.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, decks)
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
			ID:               dbCard.ID,
			FrontContent:     dbCard.FrontContent,
			BackContent:      dbCard.BackContent,
			Interval:         dbCard.Interval,
			Target:           dbCard.Target,
			EaseFactor:       dbCard.EaseFactor,
			RepetitionsCount: dbCard.RepetitionsCount,
			LastReviewedAt:   dbCard.LastReviewedAt,
			LastReviewedNum:  dbCard.LastReviewedNum,
			CreatedAt:        dbCard.CreatedAt,
			DeckID:           dbCard.DeckID,
			Tempo:            dbCard.Tempo,
			PerfectStreak:    dbCard.PerfectStreak,
			BadStreak:        dbCard.BadStreak,
		})
	}

	for _, c := range cards {
		fmt.Println("card: ", c.FrontContent)
	}
	respondWithJSON(w, http.StatusCreated, cards)
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

func (cfg *apiConfig) newDeckHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title            string    `json:"title"`
		Description      string    `json:"description"`
		UserID           uuid.UUID `json:"user_id"`
		TempoIntervalUp  int32     `json:"tempo_interval_up"`
		TempoIntervalDn  int32     `json:"tempo_interval_dn"`
		PerfectThreshold int32     `json:"perfect_threshold"`
		BadThreshold     int32     `json:"bad_threshold"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	deckParams := database.CreateDeckParams{
		Title:            params.Title,
		Description:      params.Description,
		UserID:           params.UserID,
		TempoIntervalUp:  params.TempoIntervalUp,
		TempoIntervalDn:  params.TempoIntervalDn,
		PerfectThreshold: params.PerfectThreshold,
		BadThreshold:     params.BadThreshold,
	}

	dbDeck, err := cfg.db.CreateDeck(cfg.ctx, deckParams)
	if err != nil {
		log.Fatal(err)
	}

	deck := study.Deck{
		ID:               dbDeck.ID,
		Title:            dbDeck.Title,
		Description:      dbDeck.Description,
		CreatedAt:        dbDeck.CreatedAt,
		UserID:           dbDeck.UserID,
		TotalReviews:     dbDeck.TotalReviews,
		TempoIntervalUp:  dbDeck.TempoIntervalUp,
		TempoIntervalDn:  dbDeck.TempoIntervalDn,
		PerfectThreshold: dbDeck.PerfectThreshold,
		BadThreshold:     dbDeck.BadThreshold,
	}

	respondWithJSON(w, http.StatusCreated, deck)
}
