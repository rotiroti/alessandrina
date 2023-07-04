package web

import "github.com/rotiroti/alessandrina/domain"

// AppBook is the book model used by the API.
type AppBook struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Authors   string `json:"authors"`
	Publisher string `json:"publisher"`
	Pages     int    `json:"pages"`
}

// ToAppBook converts a domain.Book to an AppBook.
func ToAppBook(book domain.Book) AppBook {
	return AppBook{
		ID:        book.ID.String(),
		Title:     book.Title,
		Authors:   book.Authors,
		Publisher: book.Publisher,
		Pages:     book.Pages,
	}
}

// AppNewBook is the new book model used by the API.
type AppNewBook struct {
	Title     string `json:"title"`
	Authors   string `json:"authors"`
	Publisher string `json:"publisher"`
	Pages     int    `json:"pages"`
}

// ToDomainNewBook converts an AppNewBook to a domain.NewBook.
func ToDomainNewBook(book AppNewBook) domain.NewBook {
	return domain.NewBook{
		Title:     book.Title,
		Authors:   book.Authors,
		Publisher: book.Publisher,
		Pages:     book.Pages,
	}
}
