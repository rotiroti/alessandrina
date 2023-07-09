package ddb

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
)

// DynamodbBook is the struct used to store books in DynamoDB.
type DynamodbBook struct {
	ID        string `dynamodbav:"id"`
	Title     string `dynamodbav:"title"`
	Authors   string `dynamodbav:"authors"`
	Publisher string `dynamodbav:"publisher"`
	Pages     int    `dynamodbav:"pages"`
	ISBN      string `dynamodbav:"isbn"`
}

// String returns a string representation of a DynamodbBook.
func (b DynamodbBook) String() string {
	const msg = "ID: %s\nTitle: %s\nAuthors: %s\nPublisher: %s\nPages: %d\nISBN: %s\n"

	return fmt.Sprintf(msg, b.ID, b.Title, b.Authors, b.Publisher, b.Pages, b.ISBN)
}

// ToDynamodbBook converts a domain.Book to a DynamodbBook.
func ToDynamodbBook(book domain.Book) DynamodbBook {
	return DynamodbBook{
		ID:        book.ID.String(),
		Title:     book.Title,
		Authors:   book.Authors,
		Publisher: book.Publisher,
		Pages:     book.Pages,
		ISBN:      book.ISBN,
	}
}

// ToDomainBook converts a DynamoDBBook to a domain.Book.
func ToDomainBook(book DynamodbBook) domain.Book {
	return domain.Book{
		ID:        uuid.MustParse(book.ID),
		Title:     book.Title,
		Authors:   book.Authors,
		Publisher: book.Publisher,
		Pages:     book.Pages,
		ISBN:      book.ISBN,
	}
}
