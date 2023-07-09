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
//
//go:generate mockery --name Storer
type Storer interface {
	Save(ctx context.Context, book Book) error
	FindOne(ctx context.Context, bookID uuid.UUID) (Book, error)
}

// Service is the domain service for Book.
type Service struct {
	storer    Storer
	generator UUIDGenerator
}

// NewService returns a new instance of Service.
func NewService(storer Storer) *Service {
	return NewServiceWithGenerator(storer, uuid.New)
}

// NewServiceWithGenerator returns a new instance of Service with a custom UUIDGenerator.
func NewServiceWithGenerator(storer Storer, generator UUIDGenerator) *Service {
	return &Service{
		storer:    storer,
		generator: generator,
	}
}

// Save adds a new book into a storage.
func (s *Service) Save(ctx context.Context, nb NewBook) (Book, error) {
	book := Book{
		ID:        s.generator(),
		Title:     nb.Title,
		Authors:   nb.Authors,
		Publisher: nb.Publisher,
		Pages:     nb.Pages,
		ISBN:      nb.ISBN,
	}

	if err := s.storer.Save(ctx, book); err != nil {
		return Book{}, fmt.Errorf("save: %w", err)
	}

	return book, nil
}

// FindOne returns a book from a storage by using bookID as primary key.
func (s *Service) FindOne(ctx context.Context, bookID uuid.UUID) (Book, error) {
	book, err := s.storer.FindOne(ctx, bookID)
	if err != nil {
		return Book{}, fmt.Errorf("findone: bookID[%s]: %w", bookID, err)
	}

	return book, nil
}
