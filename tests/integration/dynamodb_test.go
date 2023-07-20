package integration

import (
	"context"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/rotiroti/alessandrina/sys/database/ddb"
	"github.com/stretchr/testify/require"
)

func skipLocalstack(t *testing.T) {
	t.Helper()
	if os.Getenv("LOCALSTACK") == "" {
		t.Skip("skipping integration tests, set environment variable LOCALSTACK")
	}
}

func setupDB(t *testing.T) (*ddb.Store, uuid.UUID) {
	t.Helper()

	conf := ddb.Config{
		TableName: os.Getenv("TABLE_NAME"),
		Endpoint:  os.Getenv("AWS_ENDPOINT"),
		ClientLog: os.Getenv("AWS_CLIENT_DEBUG"),
	}

	store, err := ddb.NewStore(context.Background(), conf)
	require.NoError(t, err)

	bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1cabd")

	return store, bookID
}

func TestDynamodbFlow(t *testing.T) {
	// Skip the DynamoDB test if the LOCALSTACK environment variable is not set
	skipLocalstack(t)

	ctx := context.Background()
	store, expectedBookID := setupDB(t)
	expectedBook := domain.Book{
		ID:        expectedBookID,
		Title:     "The Lord of the Rings",
		Authors:   "J.R.R. Tolkien",
		Pages:     1178,
		Publisher: "George Allen & Unwin",
		ISBN:      "978-0-261-10235-4",
	}

	t.Run("Save", func(t *testing.T) {
		err := store.Save(ctx, expectedBook)
		require.NoError(t, err)
	})

	t.Run("SaveTableNotExist", func(t *testing.T) {
		conf := ddb.Config{TableName: "not-exist"}
		store, err := ddb.NewStore(ctx, conf)

		require.NoError(t, err)

		err = store.Save(ctx, domain.Book{})
		require.Error(t, err)
	})

	t.Run("FindOne", func(t *testing.T) {
		err := store.Save(ctx, expectedBook)

		require.NoError(t, err)

		book, err := store.FindOne(ctx, expectedBookID)

		require.NoError(t, err)
		require.Equal(t, expectedBook, book)
	})

	t.Run("FindOneBookNotFound", func(t *testing.T) {
		_, err := store.FindOne(ctx, uuid.MustParse("01234567-0123-0123-0123-0123456789ab"))

		require.Error(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := store.Save(ctx, expectedBook)

		require.NoError(t, err)

		err = store.Delete(ctx, expectedBookID)

		require.NoError(t, err)
	})

	t.Run("DeleteBookNotFound", func(t *testing.T) {
		err := store.Delete(ctx, uuid.MustParse("01234567-0123-0123-0123-0123456789ab"))

		require.Error(t, err)
	})

	t.Run("FindAllBooks", func(t *testing.T) {
		newBooksLen := 3
		newBooks := make([]domain.Book, newBooksLen)

		for i := 0; i < newBooksLen; i++ {
			book := domain.Book{
				ID:        uuid.New(),
				Title:     gofakeit.BookTitle(),
				Authors:   gofakeit.BookAuthor(),
				Publisher: gofakeit.Company(),
				Pages:     gofakeit.Number(100, 1200),
			}

			// Save a new book in the database
			err := store.Save(ctx, book)
			require.NoError(t, err)

			newBooks[i] = book
		}

		books, err := store.FindAll(ctx)

		// Assert the expected output
		require.NoError(t, err)
		require.LessOrEqual(t, newBooksLen, len(books))
	})
}
