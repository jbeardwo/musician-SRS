package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jbeardwo/musician-SRS/internal/study"
)

type clientConfig struct {
	baseURL       string
	currentUserID uuid.UUID
	decks         []study.Deck
}

func main() {
	clientCfg := clientConfig{
		baseURL: "http://localhost:8080",
	}

	clientCfg.clientLogin()
	clientCfg.clientGetDecks()

	for {
		fmt.Println()
		fmt.Println("Enter a command. 'help' for instructions")
		command := GetInput()
		switch command[0] {
		case "help":
			printHelp()
		case "health":
			clientCfg.healthCheck()
		case "list":
			if len(command) == 1 {
				clientCfg.listDecks()
			} else {
				num, err := strconv.Atoi(command[1])
				if err != nil {
					fmt.Printf("invalid deck")
					break
				}
				clientCfg.listCards(&clientCfg.decks[num])
			}
		case "study":
			if len(command) != 2 {
				fmt.Println("invalid usage: study [deck number]")
				break
			}
			num, err := strconv.Atoi(command[1])
			if err != nil {
				fmt.Printf("invalid deck")
				break
			}
			clientCfg.clientStudy(num)
		case "new":
			clientCfg.clientNewDeck()
		case "add":
			if len(command) != 2 {
				fmt.Println("invalid usage: study [deck number]")
				break
			}
			num, err := strconv.Atoi(command[1])
			if err != nil {
				fmt.Printf("invalid deck")
				break
			}
			clientCfg.clientNewCard(clientCfg.decks[num])
		case "delete":
			if len(command) != 2 {
				fmt.Println("invalid usage: delete [deck number]")
				break
			}
			num, err := strconv.Atoi(command[1])
			if err != nil {
				fmt.Printf("invalid deck")
				break
			}
			clientCfg.clientDeleteDeck(clientCfg.decks[num])
		case "remove":
			if len(command) != 2 {
				fmt.Println("invalid usage: remove [deck number]")
				break
			}
			deckNum, err := strconv.Atoi(command[1])
			if err != nil {
				fmt.Printf("invalid deck")
				break
			}
			clientCfg.listCards(&clientCfg.decks[deckNum])
			fmt.Println("choose a card to remove")
			command := GetInput()
			if len(command) != 1 {
				fmt.Println("invalid usage: remove [deck number]")
				break
			}
			cardNum, err := strconv.Atoi(command[0])
			if err != nil {
				fmt.Printf("invalid card")
				break
			}
			clientCfg.clientDeleteCard(clientCfg.decks[deckNum].Cards[cardNum])
		case "update":
			if len(command) != 2 {
				fmt.Println("invalid usage: update [deck number]")
				break
			}
			num, err := strconv.Atoi(command[1])
			if err != nil {
				fmt.Printf("invalid deck")
				break
			}
			clientCfg.clientUpdateDeck(clientCfg.decks[num])
		case "modify":
			if len(command) != 2 {
				fmt.Println("invalid usage: modify [deck number]")
				break
			}
			deckNum, err := strconv.Atoi(command[1])
			if err != nil {
				fmt.Printf("invalid deck")
				break
			}
			clientCfg.listCards(&clientCfg.decks[deckNum])
			fmt.Println("choose a card to modify")
			command := GetInput()
			if len(command) != 1 {
				fmt.Println("invalid usage: modify [deck number]")
				break
			}
			cardNum, err := strconv.Atoi(command[0])
			if err != nil {
				fmt.Printf("invalid card")
				break
			}
			clientCfg.clientModifyCard(clientCfg.decks[deckNum].Cards[cardNum])
		default:
			fmt.Println("invalid command")
		}
	}
}

type updateCardParams struct {
	FrontContent     string       `json:"front_content"`
	BackContent      string       `json:"back_content"`
	Interval         int32        `json:"interval"`
	Target           int32        `json:"target"`
	EaseFactor       float64      `json:"ease_factor"`
	RepetitionsCount int32        `json:"repetitions_count"`
	LastReviewedAt   sql.NullTime `json:"last_reviewed_at"`
	LastReviewedNum  int32        `json:"last_reviewed_num"`
	Tempo            int32        `json:"tempo"`
	PerfectStreak    int32        `json:"perfect_streak"`
	BadStreak        int32        `json:"bad_streak"`
	ID               uuid.UUID    `json:"id"`
}

func (cfg *clientConfig) clientModifyCard(c study.Card) {

	var params updateCardParams

	fmt.Println("Front:")
	params.FrontContent = strings.Join(GetInput(), " ")

	fmt.Println("Back:")
	params.BackContent = strings.Join(GetInput(), " ")

	params.Interval = c.Interval
	params.Target = c.Target
	params.EaseFactor = c.EaseFactor
	params.RepetitionsCount = c.RepetitionsCount
	params.LastReviewedAt = c.LastReviewedAt
	params.LastReviewedNum = c.LastReviewedNum
	params.Tempo = c.Tempo
	params.PerfectStreak = c.PerfectStreak
	params.BadStreak = c.BadStreak
	params.ID = c.ID

	_, err := requestUpdateCard(cfg.baseURL, params)
	if err != nil {
		fmt.Printf("problem updating deck: %s\n", err)
		return
	}

	cfg.clientGetDecks()
}

func requestUpdateCard(baseURL string, params updateCardParams) (study.Card, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return study.Card{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	fullURL := baseURL + "/api/cards"

	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewBuffer(payload))
	if err != nil {
		return study.Card{}, fmt.Errorf("error creating request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return study.Card{}, fmt.Errorf("Error updating deck: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return study.Card{}, fmt.Errorf("Error updating deck, status: %d\n", resp.StatusCode)
	}

	var updatedCard study.Card
	err = json.NewDecoder(resp.Body).Decode(&updatedCard)
	if err != nil {
		return study.Card{}, fmt.Errorf("error decoding deck%v\n", err)
	}

	return updatedCard, nil
}

func (cfg *clientConfig) clientDeleteDeck(d study.Deck) {
	fullURL := fmt.Sprintf("%s/api/decks/%s", cfg.baseURL, d.ID)

	req, err := http.NewRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error deleting deck: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("Failed to delete deck: status %d\n", resp.StatusCode)
		return
	}

	cfg.clientGetDecks()
}

func (cfg *clientConfig) clientDeleteCard(c study.Card) {

	fullURL := fmt.Sprintf("%s/api/cards/%s", cfg.baseURL, c.ID)

	req, err := http.NewRequest(http.MethodDelete, fullURL, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error deleting card: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Printf("Failed to delete card: status %d\n", resp.StatusCode)
		return
	}

}

func (cfg *clientConfig) clientStudy(num int) {
	if num > len(cfg.decks)-1 {
		fmt.Printf("invalid deck")
		return
	}

	cfg.clientGetCards(&cfg.decks[num])

	fmt.Println("How many cards?")

	n, err := strconv.Atoi(GetInput()[0])
	if err != nil {
		fmt.Println("Invalid Input")
		return
	}
	cfg.decks[num].ReviewDeck(n)
}

func GetInput() []string {
	fmt.Print("> ")
	scanner := bufio.NewScanner(os.Stdin)
	scanned := scanner.Scan()
	if !scanned {
		return nil
	}
	line := scanner.Text()
	line = strings.TrimSpace(line)
	return strings.Fields(line)
}

func (cfg *clientConfig) healthCheck() {
	fmt.Println("Test 1: Health Check")

	resp, err := http.Get(cfg.baseURL + "/api/healthz")
	if err != nil {
		fmt.Printf("✗ Health check failed to connect: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Health check passed")
	} else {
		fmt.Printf("Health check failed (HTTP %d)\n", resp.StatusCode)
	}
}

func (cfg *clientConfig) clientLogin() {
	fmt.Print("Enter email: ")
	var email string
	_, err := fmt.Scan(&email)
	if err != nil {
		fmt.Println("Invalid input!")
		return
	}

	payload := struct {
		Email string `json:"email"`
	}{
		Email: email,
	}

	dat, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error encoding JSON: %s\n", err)
		return
	}

	resp, err := http.Post(cfg.baseURL+"/api/login", "application/json", bytes.NewBuffer(dat))
	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Login failed: status code %d\n", resp.StatusCode)
		return
	}

	type loginResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	resBody := loginResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&resBody); err != nil {
		fmt.Printf("Error decoding response: %s\n", err)
		return
	}

	fmt.Printf("Success! Logged in as: %s (ID: %s)\n", resBody.Email, resBody.ID)
	cfg.currentUserID = resBody.ID
}

func (cfg *clientConfig) clientGetDecks() {
	fullURL := fmt.Sprintf("%s/api/decks?user_id=%s", cfg.baseURL, cfg.currentUserID)

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Error fetching decks: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get decks: status %d\n", resp.StatusCode)
		return
	}
	var decks []study.Deck

	if err := json.NewDecoder(resp.Body).Decode(&decks); err != nil {
		fmt.Printf("Error decoding decks: %v\n", err)
		return
	}

	cfg.decks = decks
}

func (cfg *clientConfig) clientGetCards(d *study.Deck) {
	fullURL := fmt.Sprintf("%s/api/cards?deck_id=%s", cfg.baseURL, d.ID)

	resp, err := http.Get(fullURL)
	if err != nil {
		fmt.Printf("Error fetching cards: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to get cards: status %d\n", resp.StatusCode)
		return
	}
	var cards []study.Card

	if err := json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		fmt.Printf("Error decoding cards: %v\n", err)
		return
	}
	d.Cards = cards
}

func (cfg *clientConfig) listDecks() {
	for i, deck := range cfg.decks {
		fmt.Println(i, " ", deck.Title)
	}
}

func (cfg *clientConfig) listCards(d *study.Deck) {
	cfg.clientGetCards(d)
	for i, card := range d.Cards {
		fmt.Println(i, " ", card.FrontContent, card.BackContent)
	}
}

type newDeckParams struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	UserID           uuid.UUID `json:"user_id"`
	TempoIntervalUp  int32     `json:"tempo_interval_up"`
	TempoIntervalDn  int32     `json:"tempo_interval_dn"`
	PerfectThreshold int32     `json:"perfect_threshold"`
	BadThreshold     int32     `json:"bad_threshold"`
}

func (cfg *clientConfig) clientNewDeck() {
	var params newDeckParams

	fmt.Println("Title:")
	params.Title = strings.Join(GetInput(), " ")

	fmt.Println("Description:")
	params.Description = strings.Join(GetInput(), " ")

	params.UserID = cfg.currentUserID

	fmt.Println("TempoIntervalUp:")
	input, err := strconv.Atoi(GetInput()[0])
	if err != nil {
		fmt.Println("Invalid Input")
		return
	}
	params.TempoIntervalUp = int32(input)

	fmt.Println("TempoIntervalDn:")
	input, err = strconv.Atoi(GetInput()[0])
	if err != nil {
		fmt.Println("Invalid Input")
		return
	}
	params.TempoIntervalDn = int32(input)

	fmt.Println("PerfectThreshold:")
	input, err = strconv.Atoi(GetInput()[0])
	if err != nil {
		fmt.Println("Invalid Input")
		return
	}
	params.PerfectThreshold = int32(input)

	fmt.Println("BadThreshold")
	input, err = strconv.Atoi(GetInput()[0])
	if err != nil {
		fmt.Println("Invalid Input")
		return
	}
	params.BadThreshold = int32(input)

	newDeck, err := requestNewDeck(cfg.baseURL, params)
	if err != nil {
		fmt.Printf("problem creating deck: %s\n", err)
	}

	cfg.decks = append(cfg.decks, newDeck)
}

func requestNewDeck(baseURL string, params newDeckParams) (study.Deck, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return study.Deck{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	fullURL := baseURL + "/api/decks"

	resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Error creating deck:%v\n", err)
		return study.Deck{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Error creating Deck: %d\n", resp.StatusCode)
	}

	var deck study.Deck

	if err := json.NewDecoder(resp.Body).Decode(&deck); err != nil {
		fmt.Printf("Error decoding deck: %v\n", err)
		return study.Deck{}, err
	}

	return deck, nil
}

type updateDeckParams struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	TotalReviews     int32     `json:"total_reviews"`
	TempoIntervalUp  int32     `json:"tempo_interval_up"`
	TempoIntervalDn  int32     `json:"tempo_interval_dn"`
	PerfectThreshold int32     `json:"perfect_threshold"`
	BadThreshold     int32     `json:"bad_threshold"`
	ID               uuid.UUID `json:"id"`
}

func (cfg *clientConfig) clientUpdateDeck(d study.Deck) {

	var params updateDeckParams

	fmt.Println("Title:")
	params.Title = strings.Join(GetInput(), " ")

	fmt.Println("Description:")
	params.Description = strings.Join(GetInput(), " ")

	fmt.Println("TempoIntervalUp:")
	input, err := strconv.Atoi(GetInput()[0])
	if err != nil {
		fmt.Println("Invalid Input")
		return
	}
	params.TempoIntervalUp = int32(input)

	fmt.Println("TempoIntervalDn:")
	input, err = strconv.Atoi(GetInput()[0])
	if err != nil {
		fmt.Println("Invalid Input")
		return
	}
	params.TempoIntervalDn = int32(input)

	fmt.Println("PerfectThreshold:")
	input, err = strconv.Atoi(GetInput()[0])
	if err != nil {
		fmt.Println("Invalid Input")
		return
	}
	params.PerfectThreshold = int32(input)

	fmt.Println("BadThreshold")
	input, err = strconv.Atoi(GetInput()[0])
	if err != nil {
		fmt.Println("Invalid Input")
		return
	}
	params.BadThreshold = int32(input)

	params.ID = d.ID
	params.TotalReviews = d.TotalReviews

	_, err = requestUpdateDeck(cfg.baseURL, params)
	if err != nil {
		fmt.Printf("problem updating deck: %s\n", err)
		return
	}

	cfg.clientGetDecks()

}

func requestUpdateDeck(baseURL string, params updateDeckParams) (study.Deck, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return study.Deck{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	fullURL := baseURL + "/api/decks"

	req, err := http.NewRequest(http.MethodPut, fullURL, bytes.NewBuffer(payload))
	if err != nil {
		return study.Deck{}, fmt.Errorf("error creating request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return study.Deck{}, fmt.Errorf("Error updating deck: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return study.Deck{}, fmt.Errorf("Error updating deck, status: %d\n", resp.StatusCode)
	}

	var updatedDeck study.Deck
	err = json.NewDecoder(resp.Body).Decode(&updatedDeck)
	if err != nil {
		return study.Deck{}, fmt.Errorf("error decoding deck%v\n", err)
	}

	return updatedDeck, nil
}

type newCardParams struct {
	FrontContent string    `json:"front_content"`
	BackContent  string    `json:"back_content"`
	Target       int32     `json:"target"`
	DeckID       uuid.UUID `json:"deck_id"`
	Tempo        int32     `json:"tempo"`
}

func (cfg *clientConfig) clientNewCard(d study.Deck) {
	var params newCardParams

	fmt.Println("FrontContent:")
	params.FrontContent = strings.Join(GetInput(), " ")

	fmt.Println("Back:")
	params.BackContent = strings.Join(GetInput(), " ")

	params.Target = d.TotalReviews

	params.DeckID = d.ID

	fmt.Println("Tempo:")
	input, err := strconv.Atoi(GetInput()[0])
	if err != nil {
		fmt.Println("Invalid Input")
		return
	}
	params.Tempo = int32(input)

	newCard, err := requestNewCard(cfg.baseURL, params)
	if err != nil {
		fmt.Printf("problem creating card: %s", err)
	}

	fmt.Printf("created card: %s\n", newCard.FrontContent)
}

func requestNewCard(baseURL string, params newCardParams) (study.Card, error) {
	payload, err := json.Marshal(params)
	if err != nil {
		return study.Card{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	fullURL := baseURL + "/api/cards"

	resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Error creating card:%v\n", err)
		return study.Card{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Error creating card: %d\n", resp.StatusCode)
	}

	var card study.Card

	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		fmt.Printf("Error decoding card: %v\n", err)
		return study.Card{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	return card, nil
}

func printHelp() {
	fmt.Println()
	fmt.Println("Commands and usage: ")
	fmt.Println()
	fmt.Println("help: show this menu")
	fmt.Println("health: check server connection")
	fmt.Println("list: list your decks alongside their deck #")
	fmt.Println("list [deck #]: list your cards within specified deck")
	fmt.Println("study [deck #]: begin reviewing specified deck")
	fmt.Println("new: creates new deck")
	fmt.Println("add [deck #]: creates a new card and adds it to the specified deck")
	fmt.Println("delete [deck #]: deletes specified deck")
	fmt.Println("update [deck #]: updates specified deck")
	fmt.Println("remove [deck #]: removes single card from specified deck")
	fmt.Println("modify [deck #]: modify a single card form specified deck")
}
