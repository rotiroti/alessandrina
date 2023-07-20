package ddb_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/rotiroti/alessandrina/sys/database/ddb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (*ddb.Store, *ddb.MockDynamodbAPI, string, uuid.UUID) {
	table := "test-table"
	client := ddb.NewMockDynamodbAPI(t)
	store := ddb.NewStore(table, client)
	bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")

	return store, client, table, bookID
}

func TestStore(t *testing.T) {
	store, client, table, bookID := setup(t)
	ctx := context.Background()
	expectedBook := domain.Book{
		ID:        bookID,
		Title:     "The Lord of the Rings",
		Authors:   "J.R.R. Tolkien",
		Pages:     1178,
		Publisher: "George Allen & Unwin",
		ISBN:      "978-0-261-10235-4",
	}
	expectedPutItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(table),
	}
	expectedKey := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: bookID.String()},
	}
	expectedGetItemInput := &dynamodb.GetItemInput{
		Key:       expectedKey,
		TableName: aws.String(table),
	}
	expectedDeleteInput := &dynamodb.DeleteItemInput{
		Key:                 expectedKey,
		TableName:           aws.String(table),
		ConditionExpression: aws.String("attribute_exists(id)"),
	}
	expectedScanInput := dynamodb.ScanInput{
		TableName: aws.String(table),
		Limit:     aws.Int32(ddb.DefaultTableScanLimit),
	}

	t.Run("Save", func(t *testing.T) {
		saveBookItem, err := attributevalue.MarshalMap(ddb.ToDynamodbBook(expectedBook))
		require.NoError(t, err)
		expectedPutItemInput.Item = saveBookItem
		client.EXPECT().PutItem(ctx, expectedPutItemInput).Return(nil, nil).Once()
		err = store.Save(ctx, expectedBook)
		require.NoError(t, err)
		client.AssertExpectations(t)
	})

	t.Run("SaveFail", func(t *testing.T) {
		saveBookItem, err := attributevalue.MarshalMap(ddb.ToDynamodbBook(expectedBook))
		require.NoError(t, err)
		expectedPutItemInput.Item = saveBookItem
		client.EXPECT().PutItem(ctx, expectedPutItemInput).Return(&dynamodb.PutItemOutput{}, assert.AnError).Once()
		err = store.Save(ctx, expectedBook)
		require.Error(t, err)
		client.AssertExpectations(t)
	})

	t.Run("FindAll", func(t *testing.T) {
		expectedScanOutput := &dynamodb.ScanOutput{
			Items: []map[string]types.AttributeValue{
				{
					"id": &types.AttributeValueMemberS{Value: bookID.String()},
				},
			},
		}
		expectedBooks := []domain.Book{
			{
				ID: uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812"),
			},
		}
		client.EXPECT().Scan(ctx, &expectedScanInput).Return(expectedScanOutput, nil).Once()
		books, err := store.FindAll(ctx)
		require.NoError(t, err)
		require.Equal(t, expectedBooks, books)
		client.AssertExpectations(t)
	})

	t.Run("FindAllFail", func(t *testing.T) {
		client.EXPECT().Scan(ctx, &expectedScanInput).Return(&dynamodb.ScanOutput{}, assert.AnError).Once()
		_, err := store.FindAll(ctx)
		require.Error(t, err)
		client.AssertExpectations(t)
	})

	t.Run("FindOne", func(t *testing.T) {
		getItemOutput, err := attributevalue.MarshalMap(ddb.ToDynamodbBook(expectedBook))
		require.NoError(t, err)
		client.EXPECT().GetItem(ctx, expectedGetItemInput).Return(
			&dynamodb.GetItemOutput{
				Item: getItemOutput,
			},
			nil,
		).Once()
		foundBook, err := store.FindOne(ctx, bookID)
		require.NoError(t, err)
		assert.Equal(t, expectedBook, foundBook)
		client.AssertExpectations(t)
	})

	t.Run("FindOneFail", func(t *testing.T) {
		client.EXPECT().GetItem(ctx, expectedGetItemInput).Return(&dynamodb.GetItemOutput{}, assert.AnError).Once()
		_, err := store.FindOne(ctx, bookID)
		require.Error(t, err)
		client.AssertExpectations(t)
	})

	t.Run("FindOneItemNotFound", func(t *testing.T) {
		client.EXPECT().GetItem(ctx, expectedGetItemInput).Return(&dynamodb.GetItemOutput{Item: map[string]types.AttributeValue{}}, nil).Once()
		_, err := store.FindOne(ctx, bookID)
		require.Error(t, err)
		client.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		client.EXPECT().DeleteItem(ctx, expectedDeleteInput).Return(nil, nil).Once()
		err := store.Delete(ctx, bookID)
		require.NoError(t, err)
		client.AssertExpectations(t)
	})
	t.Run("DeleteFail", func(t *testing.T) {
		client.EXPECT().DeleteItem(ctx, expectedDeleteInput).Return(&dynamodb.DeleteItemOutput{}, assert.AnError).Once()
		err := store.Delete(ctx, bookID)
		require.Error(t, err)
		client.AssertExpectations(t)
	})
}
