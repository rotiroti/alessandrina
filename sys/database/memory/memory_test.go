package memory_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/rotiroti/alessandrina/sys/database/memory"
	"github.com/stretchr/testify/require"
)

func TestMemoryStore(t *testing.T) {
	t.Parallel()

	book := domain.Book{
		ID:        uuid.New(),
		Title:     "The Go Programming Language",
		Authors:   "Alan A. A. Donovan, Brian W. Kernighan",
		Publisher: "Addison-Wesley Professional",
		Pages:     400,
	}

	t.Run("should save a new book", func(t *testing.T) {
		t.Parallel()
		store := memory.NewStore()
		err := store.Save(context.Background(), book)
		require.NoError(t, err)
	})

	t.Run("should not save an existing book", func(t *testing.T) {
		t.Parallel()
		store := memory.NewStore()
		err := store.Save(context.Background(), book)
		require.NoError(t, err)
		err2 := store.Save(context.Background(), book)
		require.ErrorIs(t, err2, domain.ErrAlreadyExists)
	})
}
