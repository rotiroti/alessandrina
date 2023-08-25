package ddb_test

import (
	"context"
	"os"
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

func TestNewStore(t *testing.T) {
	ctx := context.Background()
	t.Run("EmptyTableName", func(t *testing.T) {
		store, err := ddb.NewStore(ctx, "")

		require.Error(t, err)
		require.Nil(t, store)
	})

	t.Run("InvalidDefaultAWSConfig", func(t *testing.T) {
		os.Setenv("AWS_ENABLE_ENDPOINT_DISCOVERY", "foo-bar")
		store, err := ddb.NewStore(ctx, "test-table")

		require.Error(t, err)
		require.Nil(t, store)

		os.Unsetenv("AWS_ENABLE_ENDPOINT_DISCOVERY")
	})

	t.Run("WithInvalidLocalStack", func(t *testing.T) {
		os.Setenv("AWS_ENABLE_ENDPOINT_DISCOVERY", "foo-bar")

		store, err := ddb.NewStore(ctx, "test-table", ddb.WithLocalStack())

		require.Error(t, err)
		require.Nil(t, store)

		os.Unsetenv("AWS_ENABLE_ENDPOINT_DISCOVERY")
	})

	t.Run("WithInvalidClientLog", func(t *testing.T) {
		os.Setenv("AWS_ENABLE_ENDPOINT_DISCOVERY", "foo-bar")

		store, err := ddb.NewStore(ctx, "test-table", ddb.WithClientLog())

		require.Error(t, err)
		require.Nil(t, store)

		os.Unsetenv("AWS_ENABLE_ENDPOINT_DISCOVERY")
	})

	t.Run("OK", func(t *testing.T) {
		store, err := ddb.NewStore(ctx, "test-table")

		require.NoError(t, err)
		require.NotNil(t, store)
	})

	t.Run("WithClient", func(t *testing.T) {
		mockClient := ddb.NewMockDynamoDBClient(t)
		store, err := ddb.NewStore(ctx, "test-table", ddb.WithClient(mockClient))

		require.NoError(t, err)
		require.NotNil(t, store)
	})

	t.Run("WithClientLog", func(t *testing.T) {
		store, err := ddb.NewStore(ctx, "test-table", ddb.WithClientLog())

		require.NoError(t, err)
		require.NotNil(t, store)
	})

	t.Run("WithLocalStack", func(t *testing.T) {
		store, err := ddb.NewStore(ctx, "test-table", ddb.WithLocalStack())

		require.NoError(t, err)
		require.NotNil(t, store)
	})
}

func TestStore(t *testing.T) {
	ctx := context.Background()
	expectedTable := "test-table"
	mockClient := ddb.NewMockDynamoDBClient(t)
	expectedBookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")
	expectedBook := domain.Book{
		ID:        expectedBookID,
		Title:     "The Lord of the Rings",
		Authors:   "J.R.R. Tolkien",
		Pages:     1178,
		Publisher: "George Allen & Unwin",
		ISBN:      "978-0-261-10235-4",
	}
	expectedPutItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(expectedTable),
	}
	expectedKey := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: expectedBookID.String()},
	}
	expectedGetItemInput := &dynamodb.GetItemInput{
		Key:       expectedKey,
		TableName: aws.String(expectedTable),
	}
	expectedDeleteInput := &dynamodb.DeleteItemInput{
		Key:                 expectedKey,
		TableName:           aws.String(expectedTable),
		ConditionExpression: aws.String("attribute_exists(id)"),
	}
	expectedScanInput := dynamodb.ScanInput{
		TableName: aws.String(expectedTable),
		Limit:     aws.Int32(ddb.DefaultTableScanLimit),
	}

	t.Run("Save", func(t *testing.T) {
		saveBookItem, err := attributevalue.MarshalMap(ddb.ToDynamodbBook(expectedBook))
		require.NoError(t, err)
		expectedPutItemInput.Item = saveBookItem
		mockClient.EXPECT().PutItem(ctx, expectedPutItemInput).Return(nil, nil).Once()
		store, err := ddb.NewStore(ctx, expectedTable, ddb.WithClient(mockClient))
		require.NoError(t, err)
		err = store.Save(ctx, expectedBook)
		require.NoError(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("SaveFail", func(t *testing.T) {
		saveBookItem, err := attributevalue.MarshalMap(ddb.ToDynamodbBook(expectedBook))
		require.NoError(t, err)
		expectedPutItemInput.Item = saveBookItem
		mockClient.EXPECT().PutItem(ctx, expectedPutItemInput).Return(&dynamodb.PutItemOutput{}, assert.AnError).Once()
		store, err := ddb.NewStore(ctx, expectedTable, ddb.WithClient(mockClient))
		require.NoError(t, err)
		err = store.Save(ctx, expectedBook)
		require.Error(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("FindAll", func(t *testing.T) {
		expectedScanOutput := &dynamodb.ScanOutput{
			Items: []map[string]types.AttributeValue{
				{
					"id": &types.AttributeValueMemberS{Value: expectedBookID.String()},
				},
			},
		}
		expectedBooks := []domain.Book{
			{
				ID: uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812"),
			},
		}
		mockClient.EXPECT().Scan(ctx, &expectedScanInput).Return(expectedScanOutput, nil).Once()
		store, err := ddb.NewStore(ctx, expectedTable, ddb.WithClient(mockClient))
		require.NoError(t, err)
		books, err := store.FindAll(ctx)
		require.NoError(t, err)
		require.Equal(t, expectedBooks, books)
		mockClient.AssertExpectations(t)
	})

	t.Run("FindAllFail", func(t *testing.T) {
		mockClient.EXPECT().Scan(ctx, &expectedScanInput).Return(&dynamodb.ScanOutput{}, assert.AnError).Once()
		store, err := ddb.NewStore(ctx, expectedTable, ddb.WithClient(mockClient))
		require.NoError(t, err)
		_, err = store.FindAll(ctx)
		require.Error(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("FindOne", func(t *testing.T) {
		getItemOutput, err := attributevalue.MarshalMap(ddb.ToDynamodbBook(expectedBook))
		require.NoError(t, err)

		mockClient.EXPECT().GetItem(ctx, expectedGetItemInput).Return(
			&dynamodb.GetItemOutput{
				Item: getItemOutput,
			},
			nil,
		).Once()
		store, err := ddb.NewStore(ctx, expectedTable, ddb.WithClient(mockClient))
		require.NoError(t, err)
		foundBook, err := store.FindOne(ctx, expectedBookID)
		require.NoError(t, err)
		assert.Equal(t, expectedBook, foundBook)
		mockClient.AssertExpectations(t)
	})

	t.Run("FindOneFail", func(t *testing.T) {
		mockClient.EXPECT().GetItem(ctx, expectedGetItemInput).Return(&dynamodb.GetItemOutput{}, assert.AnError).Once()
		store, err := ddb.NewStore(ctx, expectedTable, ddb.WithClient(mockClient))
		require.NoError(t, err)
		_, err = store.FindOne(ctx, expectedBookID)
		require.Error(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("FindOneItemNotFound", func(t *testing.T) {
		mockClient.EXPECT().GetItem(ctx, expectedGetItemInput).Return(&dynamodb.GetItemOutput{Item: map[string]types.AttributeValue{}}, nil).Once()
		store, err := ddb.NewStore(ctx, expectedTable, ddb.WithClient(mockClient))
		require.NoError(t, err)
		_, err = store.FindOne(ctx, expectedBookID)
		require.Error(t, err)
		mockClient.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		mockClient.EXPECT().DeleteItem(ctx, expectedDeleteInput).Return(nil, nil).Once()
		store, err := ddb.NewStore(ctx, expectedTable, ddb.WithClient(mockClient))
		require.NoError(t, err)
		err = store.Delete(ctx, expectedBookID)
		require.NoError(t, err)
		mockClient.AssertExpectations(t)
	})
	t.Run("DeleteFail", func(t *testing.T) {
		mockClient.EXPECT().DeleteItem(ctx, expectedDeleteInput).Return(&dynamodb.DeleteItemOutput{}, assert.AnError).Once()
		store, err := ddb.NewStore(ctx, expectedTable, ddb.WithClient(mockClient))
		require.NoError(t, err)
		err = store.Delete(ctx, expectedBookID)
		require.Error(t, err)
		mockClient.AssertExpectations(t)
	})
}
