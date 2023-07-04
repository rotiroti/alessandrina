package web_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/rotiroti/alessandrina/sys/database/memory"
	"github.com/rotiroti/alessandrina/web"
	"github.com/stretchr/testify/require"
)

func TestCreateBookUnableToUnmarshalBodyRequest(t *testing.T) {
	t.Parallel()

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(nil)

	// Set up the expected inputs and outputs
	ctx := context.Background()

	// Call the CreateBook method of the handler
	ret, err := handler.CreateBook(ctx, events.APIGatewayV2HTTPRequest{})

	// Assert the expected output
	require.Error(t, err)

	// Assert response status code
	require.Equal(t, http.StatusBadRequest, ret.StatusCode)
}

func TestCreateBookDuplicateID(t *testing.T) {
	t.Parallel()

	// Instantiate context
	ctx := context.Background()

	// Instantiate a memory store
	store := memory.NewStore()

	// Generate a fixed UUID for the test
	expectedID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")

	// Generate a fixed UUID for the test
	mockGenerator := func() uuid.UUID {
		return expectedID
	}

	// Insert a book with the same ID
	existingBook := domain.Book{
		ID:        expectedID,
		Title:     "The Go Programming Language",
		Authors:   "Alan A. A. Donovan, Brian W. Kernighan",
		Publisher: "Addison-Wesley Professional",
		Pages:     400,
	}

	err := store.Save(ctx, existingBook)
	require.NoError(t, err)

	// Create an instance of the Service struct with the populated store and mockGenerator
	service := domain.NewServiceWithGenerator(store, mockGenerator)

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(service)

	// Set up the expected inputs and outputs
	jsonNewBook := `{
		"title": "The Go Programming Language",
		"authors": "Alan A. A. Donovan, Brian W. Kernighan",
		"publisher": "Addison-Wesley Professional",
		"pages": 400
	}`

	// Call the CreateBook method of the handler
	ret, err := handler.CreateBook(ctx, events.APIGatewayV2HTTPRequest{
		Body: jsonNewBook,
	})

	// Assert the expected output
	require.Error(t, err)

	// Assert response status code
	require.Equal(t, http.StatusInternalServerError, ret.StatusCode)
}

func TestCreateBook(t *testing.T) {
	t.Parallel()

	// Instantiate a memory store
	store := memory.NewStore()

	// Generate a fixed UUID for the test
	expectedID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")

	// Generate a fixed UUID for the test
	mockGenerator := func() uuid.UUID {
		return expectedID
	}

	// Create an instance of the Service struct with the mockStorer
	service := domain.NewServiceWithGenerator(store, mockGenerator)

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(service)

	// Set up the expected inputs and outputs
	ctx := context.Background()
	jsonNewBook := `{
		"title": "The Go Programming Language",
		"authors": "Alan A. A. Donovan, Brian W. Kernighan",
		"publisher": "Addison-Wesley Professional",
		"pages": 400
	}`
	expectedJSONBook := `{
		"id": "ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812",
		"title": "The Go Programming Language",
		"authors": "Alan A. A. Donovan, Brian W. Kernighan",
		"publisher": "Addison-Wesley Professional",
		"pages": 400
	}`

	// Call the CreateBook method of the handler
	ret, err := handler.CreateBook(ctx, events.APIGatewayV2HTTPRequest{
		Body: jsonNewBook,
	})

	// Assert the expected output
	require.NoError(t, err)

	// Assert response
	require.Equal(t, http.StatusCreated, ret.StatusCode)
	require.JSONEq(t, expectedJSONBook, ret.Body)
}
