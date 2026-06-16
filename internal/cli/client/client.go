package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

	var command string
	for {
		fmt.Println("Enter a command:")

		_, err := fmt.Scan(&command)
		if err != nil {
			fmt.Println("Invalid input!")
		}

		switch command {
		case "health":
			clientCfg.healthCheck()
		case "study":

		}
	}
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

	fmt.Printf("✓ Success! Logged in as: %s (ID: %s)\n", resBody.Email, resBody.ID)
	cfg.currentUserID = resBody.ID
	fmt.Printf("✓ Success! updated local User (ID: %s)\n", cfg.currentUserID)
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

	for _, d := range decks {
		fmt.Printf("- %s\n", d.Title)
	}
	cfg.decks = decks
}
