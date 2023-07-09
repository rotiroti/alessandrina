package memory

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
)

// Store is a simple in-memory implementation of the Storer interface.
type Store struct {
	container map[string]domain.Book
	mu        sync.RWMutex
}

// Ensure Store implements the Storer interface.
var _ domain.Storer = (*Store)(nil)

// NewStore returns a new instance of Store.
func NewStore() *Store {
	return &Store{
		container: make(map[string]domain.Book),
	}
}

// Save adds a new book into the in-memory database.
func (s *Store) Save(_ context.Context, book domain.Book) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.container[book.ID.String()]; exists {
		return domain.ErrAlreadyExists
	}

	s.container[book.ID.String()] = book

	return nil
}

// FindOne returns a book from the in-memory database.
func (s *Store) FindOne(_ context.Context, bookID uuid.UUID) (domain.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, exists := s.container[bookID.String()]
	if !exists {
		return domain.Book{}, domain.ErrNotFound
	}

	return book, nil
}
