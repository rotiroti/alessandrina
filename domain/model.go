package domain

import "github.com/google/uuid"

// Book represents information about an individual book.
type Book struct {
	ID        uuid.UUID
	Title     string
	Authors   string
	Publisher string
	Pages     int
	ISBN      string
}

// NewBook contains information needed to create a new book.
type NewBook struct {
	Title     string
	Authors   string
	Publisher string
	Pages     int
	ISBN      string
}
