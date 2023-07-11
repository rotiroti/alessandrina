package web_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/rotiroti/alessandrina/sys/database/memory"
	"github.com/rotiroti/alessandrina/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	require.NoError(t, err)

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
		ISBN:      "978-0134190440",
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
		"pages": 400,
		"isbn": "978-0134190440"
	}`

	// Call the CreateBook method of the handler
	ret, err := handler.CreateBook(ctx, events.APIGatewayV2HTTPRequest{
		Body: jsonNewBook,
	})

	// Assert the expected output
	require.NoError(t, err)

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
		"pages": 400,
		"isbn": "978-0134190440"
	}`
	expectedJSONBook := `{
		"id": "ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812",
		"title": "The Go Programming Language",
		"authors": "Alan A. A. Donovan, Brian W. Kernighan",
		"publisher": "Addison-Wesley Professional",
		"pages": 400,
		"isbn": "978-0134190440"
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

func TestGetBookUnableToUnmarshalPathParameters(t *testing.T) {
	t.Parallel()

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(nil)

	// Set up the expected inputs and outputs
	ctx := context.Background()

	// Call the GetBook method of the handler
	ret, err := handler.GetBook(ctx, events.APIGatewayV2HTTPRequest{})

	// Assert the expected output
	require.NoError(t, err)

	// Assert the expected output
	require.Equal(t, http.StatusBadRequest, ret.StatusCode)
}

func TestGetBookInternalServerError(t *testing.T) {
	t.Parallel()

	// Instantiate a memory store
	store := memory.NewStore()

	// Instantiate a new service
	service := domain.NewService(store)

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(service)

	// Set up the expected inputs and outputs
	ctx := context.Background()
	parameterID := "ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812"

	// Call the GetBook method of the handler
	ret, err := handler.GetBook(ctx, events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"id": parameterID,
		},
	})

	// Assert the expected output
	require.NoError(t, err)

	// Assert the expected output
	require.Equal(t, http.StatusInternalServerError, ret.StatusCode)
}

func TestGetBook(t *testing.T) {
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
		ISBN:      "978-0134190440",
	}

	err := store.Save(ctx, existingBook)
	require.NoError(t, err)

	// Create an instance of the Service struct with the populated store and mockGenerator
	service := domain.NewServiceWithGenerator(store, mockGenerator)

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(service)

	// Set up the expected inputs and outputs
	expectedJSONBook := `{
		"id": "ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812",
		"title": "The Go Programming Language",
		"authors": "Alan A. A. Donovan, Brian W. Kernighan",
		"publisher": "Addison-Wesley Professional",
		"pages": 400,
		"isbn": "978-0134190440"
	}`

	// Call the CreateBook method of the handler
	ret, err := handler.GetBook(ctx, events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"id": expectedID.String(),
		},
	})

	// Assert the expected output
	require.NoError(t, err)

	// Assert response
	require.Equal(t, http.StatusOK, ret.StatusCode)
	require.JSONEq(t, expectedJSONBook, ret.Body)
}

func TestDeleteBookUnableToUnmarshalPathParameters(t *testing.T) {
	t.Parallel()

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(nil)

	// Set up the expected inputs and outputs
	ctx := context.Background()

	// Call the DeleteBook method of the handler
	ret, err := handler.DeleteBook(ctx, events.APIGatewayV2HTTPRequest{})

	// Assert the expected output
	require.NoError(t, err)

	// Assert the expected output
	require.Equal(t, http.StatusBadRequest, ret.StatusCode)
}

func TestDeleteBookInternalServerError(t *testing.T) {
	t.Parallel()

	// Instantiate a memory store
	store := memory.NewStore()

	// Instantiate a new service
	service := domain.NewService(store)

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(service)

	// Set up the expected inputs and outputs
	ctx := context.Background()
	parameterID := "ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812"

	// Call the DeleteBook method of the handler
	ret, err := handler.DeleteBook(ctx, events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"id": parameterID,
		},
	})

	// Assert the expected output
	require.NoError(t, err)

	// Assert the expected output
	require.Equal(t, http.StatusInternalServerError, ret.StatusCode)
}

func TestDelete(t *testing.T) {
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
		Title:     "The Lord of the Rings",
		Authors:   "J.R.R. Tolkien",
		Publisher: "George Allen & Unwin",
		Pages:     1178,
	}

	err := store.Save(ctx, existingBook)
	require.NoError(t, err)

	// Create an instance of the Service struct with the populated store and mockGenerator
	service := domain.NewServiceWithGenerator(store, mockGenerator)

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(service)

	// Call the DeleteBook method of the handler
	ret, err := handler.DeleteBook(ctx, events.APIGatewayV2HTTPRequest{
		PathParameters: map[string]string{
			"id": expectedID.String(),
		},
	})

	// Assert the expected output
	require.NoError(t, err)

	// Assert the expected output
	require.Equal(t, http.StatusNoContent, ret.StatusCode)
}

func TestGetBooks(t *testing.T) {
	t.Parallel()

	// Instantiate context
	ctx := context.Background()

	// Instantiate a memory store
	store := memory.NewStore()

	// Seed the store with some books
	const expectedBookCount = 5

	for i := 0; i < expectedBookCount; i++ {
		book := domain.Book{
			ID:        uuid.New(),
			Title:     gofakeit.BookTitle(),
			Authors:   gofakeit.BookAuthor(),
			Publisher: gofakeit.Company(),
			Pages:     gofakeit.Number(100, 1200),
		}

		err := store.Save(ctx, book)
		require.NoError(t, err)
	}

	// Create an instance of the Service struct with the populated store and mockGenerator
	service := domain.NewService(store)

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(service)

	// Call the GetBooks method of the handler
	ret, err := handler.GetBooks(ctx, events.APIGatewayV2HTTPRequest{})

	// Assert the expected output
	require.NoError(t, err)

	// Assert response
	require.Equal(t, http.StatusOK, ret.StatusCode)
}

type MockStorer struct {
	mock.Mock
}

func (m *MockStorer) FindAll(ctx context.Context) ([]domain.Book, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Book), args.Error(1)
}

func (m *MockStorer) FindOne(ctx context.Context, bookID uuid.UUID) (domain.Book, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(domain.Book), args.Error(1)
}

func (m *MockStorer) Save(ctx context.Context, book domain.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *MockStorer) Delete(ctx context.Context, bookID uuid.UUID) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func TestGetBooksFail(t *testing.T) {
	t.Parallel()

	// Instantiate context
	ctx := context.Background()

	// Instantiate a mock store
	store := new(MockStorer)

	// Set up the expected inputs and outputs
	store.On("FindAll", ctx).Return([]domain.Book{}, assert.AnError).Once()

	// Instantiate a new service
	service := domain.NewService(store)

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(service)

	// Call the GetBooks method of the handler
	ret, err := handler.GetBooks(ctx, events.APIGatewayV2HTTPRequest{})

	// Assert the expected output
	require.NoError(t, err)

	// Assert response
	require.Equal(t, http.StatusInternalServerError, ret.StatusCode)
}
