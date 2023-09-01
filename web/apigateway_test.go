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

func setup(t *testing.T) (uuid.UUID, func() uuid.UUID) {
	bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")
	generator := func() uuid.UUID {
		return bookID
	}

	return bookID, generator
}

func TestBadRequest(t *testing.T) {
	type testCase struct {
		name   string
		handle func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error)
	}

	handler := web.NewAPIGatewayV2Handler(nil)
	ctx := context.Background()
	testCases := []testCase{
		{name: "CreateBook", handle: handler.CreateBook},
		{name: "GetBook", handle: handler.GetBook},
		{name: "DeleteBook", handle: handler.DeleteBook},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ret, err := tc.handle(ctx, events.APIGatewayV2HTTPRequest{})

			require.NoError(t, err)
			require.Equal(t, http.StatusBadRequest, ret.StatusCode)
		})
	}
}

func TestCreateBookInvalidPayload(t *testing.T) {
	ctx := context.Background()
	handler := web.NewAPIGatewayV2Handler(nil)
	ret, err := handler.CreateBook(ctx, events.APIGatewayV2HTTPRequest{
		Body: "{}",
	})

	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, ret.StatusCode)
}

func TestHandler(t *testing.T) {
	ctx := context.Background()
	expectedID, generator := setup(t)
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
	existingBook := domain.Book{
		ID:        expectedID,
		Title:     "The Go Programming Language",
		Authors:   "Alan A. A. Donovan, Brian W. Kernighan",
		Publisher: "Addison-Wesley Professional",
		Pages:     400,
		ISBN:      "978-0134190440",
	}

	t.Run("CreateBookDuplicateID", func(t *testing.T) {
		store := memory.NewStore()
		err := store.Save(ctx, existingBook)

		require.NoError(t, err)

		bookCore := domain.NewBookCoreWithGenerator(store, generator)
		handler := web.NewAPIGatewayV2Handler(bookCore)
		ret, err := handler.CreateBook(ctx, events.APIGatewayV2HTTPRequest{
			Body: jsonNewBook,
		})

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, ret.StatusCode)
	})

	t.Run("CreateBook", func(t *testing.T) {
		store := memory.NewStore()
		bookCore := domain.NewBookCoreWithGenerator(store, generator)
		handler := web.NewAPIGatewayV2Handler(bookCore)
		ret, err := handler.CreateBook(ctx, events.APIGatewayV2HTTPRequest{
			Body: jsonNewBook,
		})

		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, ret.StatusCode)
		require.JSONEq(t, expectedJSONBook, ret.Body)
	})

	t.Run("GetBookNotFound", func(t *testing.T) {
		store := memory.NewStore()
		bookCore := domain.NewBookCore(store)
		handler := web.NewAPIGatewayV2Handler(bookCore)
		parameterID := "ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812"
		ret, err := handler.GetBook(ctx, events.APIGatewayV2HTTPRequest{
			PathParameters: map[string]string{
				"id": parameterID,
			},
		})

		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, ret.StatusCode)
	})

	t.Run("GetBookInternalServerError", func(t *testing.T) {
		store := new(MockStorer)
		store.On("FindOne", ctx, expectedID).Return(domain.Book{}, assert.AnError).Once()
		bookCore := domain.NewBookCore(store)
		handler := web.NewAPIGatewayV2Handler(bookCore)
		ret, err := handler.GetBook(ctx, events.APIGatewayV2HTTPRequest{
			PathParameters: map[string]string{
				"id": expectedID.String(),
			},
		})

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, ret.StatusCode)
	})

	t.Run("GetBook", func(t *testing.T) {
		store := memory.NewStore()
		err := store.Save(ctx, existingBook)

		require.NoError(t, err)

		bookCore := domain.NewBookCoreWithGenerator(store, generator)
		handler := web.NewAPIGatewayV2Handler(bookCore)
		ret, err := handler.GetBook(ctx, events.APIGatewayV2HTTPRequest{
			PathParameters: map[string]string{
				"id": expectedID.String(),
			},
		})

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, ret.StatusCode)
		require.JSONEq(t, expectedJSONBook, ret.Body)
	})

	t.Run("DeleteBookNotFound", func(t *testing.T) {
		store := memory.NewStore()
		bookCore := domain.NewBookCore(store)
		handler := web.NewAPIGatewayV2Handler(bookCore)
		parameterID := "ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812"
		ret, err := handler.DeleteBook(ctx, events.APIGatewayV2HTTPRequest{
			PathParameters: map[string]string{
				"id": parameterID,
			},
		})

		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, ret.StatusCode)
	})

	t.Run("DeleteBookInternalServerError", func(t *testing.T) {
		store := new(MockStorer)
		store.On("Delete", ctx, expectedID).Return(assert.AnError).Once()
		bookCore := domain.NewBookCore(store)
		handler := web.NewAPIGatewayV2Handler(bookCore)
		ret, err := handler.DeleteBook(ctx, events.APIGatewayV2HTTPRequest{
			PathParameters: map[string]string{
				"id": expectedID.String(),
			},
		})

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, ret.StatusCode)
	})

	t.Run("DeleteBook", func(t *testing.T) {
		store := memory.NewStore()
		err := store.Save(ctx, existingBook)

		require.NoError(t, err)

		bookCore := domain.NewBookCoreWithGenerator(store, generator)
		handler := web.NewAPIGatewayV2Handler(bookCore)
		ret, err := handler.DeleteBook(ctx, events.APIGatewayV2HTTPRequest{
			PathParameters: map[string]string{
				"id": expectedID.String(),
			},
		})

		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, ret.StatusCode)
	})

	t.Run("GetBooksFail", func(t *testing.T) {
		store := new(MockStorer)
		store.On("FindAll", ctx).Return([]domain.Book{}, assert.AnError).Once()
		bookCore := domain.NewBookCore(store)
		handler := web.NewAPIGatewayV2Handler(bookCore)
		ret, err := handler.GetBooks(ctx, events.APIGatewayV2HTTPRequest{})

		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, ret.StatusCode)
	})

	t.Run("GetBooks", func(t *testing.T) {
		const expectedBookCount = 5

		store := memory.NewStore()

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

		bookCore := domain.NewBookCore(store)
		handler := web.NewAPIGatewayV2Handler(bookCore)
		ret, err := handler.GetBooks(ctx, events.APIGatewayV2HTTPRequest{})

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, ret.StatusCode)
	})
}
