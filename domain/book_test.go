package domain_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*domain.MockStorer, uuid.UUID, func() uuid.UUID) {
	storer := domain.NewMockStorer(t)
	bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")
	generator := func() uuid.UUID {
		return bookID
	}

	return storer, bookID, generator
}

func TestBookCore(t *testing.T) {
	storer, expectedID, generator := setup(t)
	ctx := context.Background()
	core := domain.NewBookCore(storer)
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

	t.Run("Save", func(t *testing.T) {
		coreWithGenerator := domain.NewBookCoreWithGenerator(storer, generator)
		storer.EXPECT().Save(ctx, expectedBook).Return(nil).Once()
		createdBook, err := coreWithGenerator.Save(ctx, newBook)
		assert.NoError(t, err)
		assert.Equal(t, expectedBook, createdBook)
		storer.AssertExpectations(t)
	})

	t.Run("SaveFail", func(t *testing.T) {
		coreWithGenerator := domain.NewBookCoreWithGenerator(storer, generator)
		storer.EXPECT().Save(ctx, expectedBook).Return(assert.AnError).Once()
		createdBook, err := coreWithGenerator.Save(ctx, newBook)
		assert.Error(t, err)
		assert.Equal(t, domain.Book{}, createdBook)
		storer.AssertExpectations(t)
	})

	t.Run("FindAll", func(t *testing.T) {
		expectedBooks := []domain.Book{
			{
				ID: uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812"),
			},
			{
				ID: uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c813"),
			},
			{
				ID: uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c814"),
			},
		}

		storer.EXPECT().FindAll(ctx).Return(expectedBooks, nil).Once()
		foundBooks, err := core.FindAll(ctx)
		assert.NoError(t, err)
		assert.Equal(t, len(expectedBooks), len(foundBooks))
		assert.Equal(t, expectedBooks, foundBooks)
		storer.AssertExpectations(t)
	})

	t.Run("FindAllFail", func(t *testing.T) {
		expectedBooks := []domain.Book{}
		storer.EXPECT().FindAll(ctx).Return(expectedBooks, assert.AnError).Once()
		foundBooks, err := core.FindAll(ctx)
		assert.Error(t, err)
		assert.Equal(t, expectedBooks, foundBooks)
		storer.AssertExpectations(t)
	})

	t.Run("FindOne", func(t *testing.T) {
		storer.EXPECT().FindOne(ctx, expectedID).Return(expectedBook, nil).Once()
		foundBook, err := core.FindOne(ctx, expectedID)
		assert.NoError(t, err)
		assert.Equal(t, expectedBook, foundBook)
		storer.AssertExpectations(t)
	})

	t.Run("FindOneFail ", func(t *testing.T) {
		ret := domain.Book{}
		storer.EXPECT().FindOne(ctx, expectedID).Return(ret, assert.AnError).Once()
		foundBook, err := core.FindOne(ctx, expectedID)
		assert.Error(t, err)
		assert.Equal(t, ret, foundBook)
		storer.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		storer.EXPECT().Delete(ctx, expectedID).Return(nil).Once()
		err := core.Delete(ctx, expectedID)
		assert.NoError(t, err)
		storer.AssertExpectations(t)
	})

	t.Run("DeleteFail", func(t *testing.T) {
		storer.EXPECT().Delete(ctx, expectedID).Return(assert.AnError).Once()
		err := core.Delete(ctx, expectedID)
		assert.Error(t, err)
		storer.AssertExpectations(t)
	})
}
