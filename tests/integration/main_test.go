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

const baseURLPath = "/books"

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

func setup() string {
	u := os.Getenv("API_URL")
	if u == "" {
		panic("API_URL environment variable is not set")
	}

	// Remove the trailing slash from the URL if it exists
	if u[len(u)-1] == '/' {
		u = u[:len(u)-1]
	}

	baseURL := fmt.Sprintf("%s%s", u, baseURLPath)

	return baseURL
}

func TestIntegrationFlow(t *testing.T) {
	t.Parallel()

	// Skip the integration test if the INTEGRATION environment variable is not set
	skipIntegration(t)

	// Setup test environment
	baseURL := setup()

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
	req, err := http.NewRequest(http.MethodPost, baseURL, bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Set any necessary headers for the request
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// Create an HTTP client and make the request
	client := &http.Client{}
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

	bookID, ok := responseBody["id"].(string)
	if !ok {
		t.Fatalf("Failed to get book ID from response: %v", err)
	}

	// Check the ID field for a valid UUID
	if _, err := uuid.Parse(bookID); err != nil {
		t.Errorf("Invalid ID format. Expected a valid UUIDv4 but got %q", bookID)
	}

	// --- GetBook scenario ---
	bookURL := fmt.Sprintf("%s/%s", baseURL, bookID)
	req, err = http.NewRequest(http.MethodGet, bookURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

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
