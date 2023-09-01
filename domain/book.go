package domain

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	// ErrNotFound is used when a specific Book is requested but does not exist.
	ErrNotFound = errors.New("book not found")

	// ErrAlreadyExists is used when a specific Book is created but already exists.
	ErrAlreadyExists = errors.New("book already exists")
)

// UUIDGenerator is a function that returns a UUID.
type UUIDGenerator func() uuid.UUID

// Storer is the interface used to interact with a storage.
type Storer interface {
	Save(ctx context.Context, book Book) error
	FindAll(ctx context.Context) ([]Book, error)
	FindOne(ctx context.Context, bookID uuid.UUID) (Book, error)
	Delete(ctx context.Context, bookID uuid.UUID) error
}

// BookCore manages the set of APIs for book access.
type BookCore struct {
	storer    Storer
	generator UUIDGenerator
}

// NewBookCore constructs a core for book API access.
func NewBookCore(storer Storer) *BookCore {
	return NewBookCoreWithGenerator(storer, uuid.New)
}

// NewBookCore constructs a core for book API access with a custom UUIDGenerator.
func NewBookCoreWithGenerator(storer Storer, generator UUIDGenerator) *BookCore {
	return &BookCore{
		storer:    storer,
		generator: generator,
	}
}

// Save inserts a new book into a storage.
func (c *BookCore) Save(ctx context.Context, nb NewBook) (Book, error) {
	book := Book{
		ID:        c.generator(),
		Title:     nb.Title,
		Authors:   nb.Authors,
		Publisher: nb.Publisher,
		Pages:     nb.Pages,
		ISBN:      nb.ISBN,
	}

	if err := c.storer.Save(ctx, book); err != nil {
		return Book{}, fmt.Errorf("domain.save failed: %w", err)
	}

	return book, nil
}

// FindAll returns all books from a storage.
func (c *BookCore) FindAll(ctx context.Context) ([]Book, error) {
	books, err := c.storer.FindAll(ctx)
	if err != nil {
		return []Book{}, fmt.Errorf("domain.findall failed: %w", err)
	}

	return books, nil
}

// FindOne returns a book from a storage by using bookID as primary key.
func (c *BookCore) FindOne(ctx context.Context, bookID uuid.UUID) (Book, error) {
	book, err := c.storer.FindOne(ctx, bookID)
	if err != nil {
		return Book{}, fmt.Errorf("domain.findone failed: %w", err)
	}

	return book, nil
}

// Delete removes a book from a storage by using bookID as primary key.
func (c *BookCore) Delete(ctx context.Context, bookID uuid.UUID) error {
	if err := c.storer.Delete(ctx, bookID); err != nil {
		return fmt.Errorf("domain.delete failed: %w", err)
	}

	return nil
}
