package domain_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	t.Parallel()

	// Create a mock instance of the Storer interface
	mockStorer := domain.NewMockStorer(t)

	// Generate a fixed UUID for the test
	expectedID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")

	// Generate a fixed UUID for the test
	mockGenerator := func() uuid.UUID {
		return expectedID
	}

	// Create an instance of the Service struct with the mockStorer
	service := domain.NewServiceWithGenerator(mockStorer, mockGenerator)

	// Set up the expected inputs and outputs
	ctx := context.TODO()
	newBook := domain.NewBook{
		Title:     "Test Book",
		Authors:   "Test Author",
		Publisher: "Test Publisher",
		Pages:     100,
		ISBN:      "Test ISBN",
	}
	expectedBook := domain.Book{
		ID:        expectedID,
		Title:     newBook.Title,
		Authors:   newBook.Authors,
		Publisher: newBook.Publisher,
		Pages:     newBook.Pages,
		ISBN:      newBook.ISBN,
	}

	// Set up the expectations for the mockStorer's Save method
	mockStorer.EXPECT().Save(ctx, expectedBook).Return(nil).Once()

	// Call the Save method of the service
	createdBook, err := service.Save(ctx, newBook)

	// Assert the expected output
	assert.NoError(t, err)
	assert.Equal(t, expectedBook, createdBook)

	// Assert that the mockStorer's expectations were met
	mockStorer.AssertExpectations(t)
}

func TestSaveFail(t *testing.T) {
	t.Parallel()

	// Create a mock instance of the Storer interface
	mockStorer := domain.NewMockStorer(t)

	// Generate a fixed UUID for the test
	expectedID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")

	// Generate a fixed UUID for the test
	mockGenerator := func() uuid.UUID {
		return expectedID
	}

	// Create an instance of the Service struct with the mockStorer
	service := domain.NewServiceWithGenerator(mockStorer, mockGenerator)

	// Set up the expected inputs and outputs
	ctx := context.TODO()
	newBook := domain.NewBook{
		Title:     "Test Book",
		Authors:   "Test Author",
		Publisher: "Test Publisher",
		Pages:     100,
		ISBN:      "Test ISBN",
	}
	expectedBook := domain.Book{
		ID:        expectedID,
		Title:     newBook.Title,
		Authors:   newBook.Authors,
		Publisher: newBook.Publisher,
		Pages:     newBook.Pages,
		ISBN:      newBook.ISBN,
	}

	// Set up the expectations for the mockStorer's Save method
	mockStorer.EXPECT().Save(ctx, expectedBook).Return(assert.AnError).Once()

	// Call the Save method of the service
	createdBook, err := service.Save(ctx, newBook)

	// Assert the expected output
	assert.Error(t, err)
	assert.Equal(t, domain.Book{}, createdBook)

	// Assert that the mockStorer's expectations were met
	mockStorer.AssertExpectations(t)
}

func TestFindOne(t *testing.T) {
	t.Parallel()

	// Create a mock instance of the Storer interface
	mockStorer := domain.NewMockStorer(t)

	// Generate a fixed UUID for the test
	expectedID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")

	// Create an instance of the Service struct with the mockStorer
	service := domain.NewService(mockStorer)

	// Set up the expected inputs and outputs
	ctx := context.TODO()
	expectedBook := domain.Book{
		ID:        expectedID,
		Title:     "Test Book",
		Authors:   "Test Authors",
		Publisher: "Test Publisher",
		ISBN:      "Test ISBN",
		Pages:     100,
	}

	// Set up the expectations for the mockStorer's FindOne method
	mockStorer.EXPECT().FindOne(ctx, expectedID).Return(expectedBook, nil).Once()

	// Call the FindOne method of the service
	foundBook, err := service.FindOne(ctx, expectedID)

	// Assert the expected output
	assert.NoError(t, err)
	assert.Equal(t, expectedBook, foundBook)

	// Assert that the mockStorer's expectations were met
	mockStorer.AssertExpectations(t)
}

func TestFindOneFail(t *testing.T) {
	t.Parallel()

	// Create a mock instance of the Storer interface
	mockStorer := domain.NewMockStorer(t)

	// Generate a fixed UUID for the test
	bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")

	// Create an instance of the Service struct with the mockStorer
	service := domain.NewService(mockStorer)

	// Set up the expected inputs and outputs
	ctx := context.TODO()

	// Set up the expectations for the mockStorer's FindOne method
	mockStorer.EXPECT().FindOne(ctx, bookID).Return(domain.Book{}, assert.AnError).Once()

	// Call the FindOne method of the service
	foundBook, err := service.FindOne(ctx, bookID)

	// Assert the expected output
	assert.Error(t, err)
	assert.Equal(t, domain.Book{}, foundBook)

	// Assert that the mockStorer's expectations were met
	mockStorer.AssertExpectations(t)
}
