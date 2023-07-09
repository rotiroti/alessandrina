package integration

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

func skipIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("skipping integration tests, set environment variable INTEGRATION")
	}
}

func generateRandomISBN() string {
	const isbnDigits = 9999999999999
	n, err := rand.Int(rand.Reader, big.NewInt(isbnDigits))
	if err != nil {
		panic(err)
	}
	return n.String()
}

func TestIntegrationFlow(t *testing.T) {
	t.Parallel()

	// Skip the integration test if the GO_RUN_INTEGRATION environment variable is not set
	skipIntegration(t)

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		t.Fatal("API_URL environment variable is not set")
	}

	// Remove the trailing slash from the URL if it exists
	if apiURL[len(apiURL)-1] == '/' {
		apiURL = apiURL[:len(apiURL)-1]
	}

	// Generate a random JSON payload for the book data
	bookData := map[string]interface{}{
		"title":     gofakeit.BookTitle(),
		"authors":   gofakeit.BookAuthor(),
		"publisher": gofakeit.Company(),
		"isbn":      generateRandomISBN(),
		"pages":     gofakeit.Number(100, 1200),
	}

	payload, err := json.Marshal(bookData)
	if err != nil {
		t.Fatalf("Failed to marshal book data: %v", err)
	}

	// --- CreateBook scenario ---
	req, err := http.NewRequest(http.MethodPost, apiURL+"/books", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set any necessary headers for the request
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and make the request
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code to be 201
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d but got %d", http.StatusCreated, resp.StatusCode)
	}

	// Parse the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Check the ID field for a valid UUID
	if _, err := uuid.Parse(responseBody["id"].(string)); err != nil {
		t.Errorf("Invalid ID format. Expected a valid UUIDv4 but got %q", responseBody["id"])
	}

	// --- GetBook scenario ---
	getBookURL := fmt.Sprintf("%s/books/%s", apiURL, responseBody["id"])

	req, err = http.NewRequest(http.MethodGet, getBookURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code to be 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}
}
