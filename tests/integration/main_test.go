package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
	isbnList := []string{
		"9780134190440",
		"978-0134190440",
		"9780321601919",
		"978-0321601919",
	}

	return isbnList[gofakeit.Number(0, len(isbnList)-1)]
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
		t.Fatalf("Expected status code %d but got %d", http.StatusCreated, resp.StatusCode)
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
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}

	// --- GetBooks scenario ---
	req, err = http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	// Check the response status code to be 200
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}

	// --- DeleteBook scenario ---
	req, err = http.NewRequest(http.MethodDelete, bookURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the response status code to be 204
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Expected status code %d but got %d", http.StatusNoContent, resp.StatusCode)
	}
}

func TestErrorResponses(t *testing.T) {
	skipIntegration(t)

	type args struct {
		method string
		url    string
		body   io.Reader
	}

	client := &http.Client{}
	baseURL := setup()
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "GetBookInvalidIdFormat",
			args: args{
				method: http.MethodGet,
				url:    fmt.Sprintf("%s/%s", baseURL, "1234"),
				body:   nil,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "CreateBookInvalidPayload",
			args: args{
				method: http.MethodPost,
				url:    baseURL,
				body:   strings.NewReader("invalid"),
			},
			want: http.StatusBadRequest,
		},
		{
			name: "CreateBookFailedValidation",
			args: args{
				method: http.MethodPost,
				url:    baseURL,
				body:   strings.NewReader(`{"title": ""}`),
			},
			want: http.StatusBadRequest,
		},
		{
			name: "DeleteBookInvalidIdFormat",
			args: args{
				method: http.MethodDelete,
				url:    fmt.Sprintf("%s/%s", baseURL, "1234"),
				body:   nil,
			},
			want: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.args.method, tt.args.url, tt.args.body)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", "application/json; charset=utf-8")
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}

			defer resp.Body.Close()

			if resp.StatusCode != tt.want {
				t.Fatalf("Expected status code %d but got %d", tt.want, resp.StatusCode)
			}
		})
	}
}
