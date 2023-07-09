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

func TestDynamoDBStore_Save(t *testing.T) {
	t.Parallel()

	tableName := "test-table"

	// Create a mock instance of the DynamoDB interface
	mockDynamoDB := ddb.NewMockDynamodbAPI(t)

	// Create a new store using the mock DynamoDB instance
	store := ddb.NewStore(tableName, mockDynamoDB)

	// Set up the expected inputs and outputs
	ctx := context.Background()
	bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")
	book := domain.Book{
		ID:        bookID,
		Title:     "The Lord of the Rings",
		Authors:   "J.R.R. Tolkien",
		Pages:     1178,
		Publisher: "George Allen & Unwin",
		ISBN:      "978-0-261-10235-4",
	}
	saveBookItem, err := attributevalue.MarshalMap(ddb.ToDynamodbBook(book))
	require.NoError(t, err)

	expectedInput := &dynamodb.PutItemInput{
		Item:      saveBookItem,
		TableName: aws.String(tableName),
	}

	// Expect a call to PutItem with the expected input and return no error
	mockDynamoDB.EXPECT().PutItem(ctx, expectedInput).Return(nil, nil).Once()

	// Call the Save method of the store
	err = store.Save(context.Background(), book)

	// Assert the expected output
	require.NoError(t, err)

	// Assert that the mockDynamoDB's expectations were met
	mockDynamoDB.AssertExpectations(t)
}

func TestDynamoDBStore_SaveFail(t *testing.T) {
	t.Parallel()

	tableName := "test-table"

	// Create a mock instance of the DynamoDB interface
	mockDynamoDB := ddb.NewMockDynamodbAPI(t)

	// Create a new store using the mock DynamoDB instance
	store := ddb.NewStore(tableName, mockDynamoDB)

	// Set up the expected inputs and outputs
	ctx := context.Background()
	bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")
	book := domain.Book{
		ID:        bookID,
		Title:     "The Lord of the Rings",
		Authors:   "J.R.R. Tolkien",
		Pages:     1178,
		Publisher: "George Allen & Unwin",
		ISBN:      "978-0-261-10235-4",
	}

	saveBookItem, err := attributevalue.MarshalMap(ddb.ToDynamodbBook(book))
	require.NoError(t, err)

	expectedInput := &dynamodb.PutItemInput{
		Item:      saveBookItem,
		TableName: aws.String(tableName),
	}

	// Expect a call to PutItem with the expected input and return no error
	mockDynamoDB.EXPECT().PutItem(ctx, expectedInput).Return(&dynamodb.PutItemOutput{}, assert.AnError).Once()

	// Call the Save method of the store
	err = store.Save(context.Background(), book)

	// Assert the expected output
	require.Error(t, err)

	// Assert that the mockDynamoDB's expectations were met
	mockDynamoDB.AssertExpectations(t)
}

func TestDynamoDBStore_FindOne(t *testing.T) {
	t.Parallel()

	tableName := "test-table"

	// Create a mock instance of the DynamoDB interface
	mockDynamoDB := ddb.NewMockDynamodbAPI(t)

	// Create a new store using the mock DynamoDB instance
	store := ddb.NewStore(tableName, mockDynamoDB)

	// Set up the expected inputs and outputs
	ctx := context.Background()
	bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")
	keyItem := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: bookID.String()},
	}
	expectedGetItemInput := &dynamodb.GetItemInput{
		Key:       keyItem,
		TableName: aws.String(tableName),
	}

	// Create a mock DynamoDB GetItemOutput
	expectedBook := domain.Book{
		ID:        bookID,
		Title:     "The Lord of the Rings",
		Authors:   "J.R.R. Tolkien",
		Pages:     1178,
		Publisher: "George Allen & Unwin",
		ISBN:      "978-0-261-10235-4",
	}

	getItemOutput, err := attributevalue.MarshalMap(ddb.ToDynamodbBook(expectedBook))
	require.NoError(t, err)

	// Expect a call to GetItem with the expected input and return the mock output
	mockDynamoDB.EXPECT().GetItem(ctx, expectedGetItemInput).Return(
		&dynamodb.GetItemOutput{
			Item: getItemOutput,
		},
		nil,
	).Once()

	// Call the FindOne method of the store
	foundBook, err := store.FindOne(ctx, bookID)

	// Assert the expected output
	require.NoError(t, err)
	assert.Equal(t, expectedBook, foundBook)

	// Assert that the mockDynamoDB's expectations were met
	mockDynamoDB.AssertExpectations(t)
}

func TestDynamoDBStore_FindOneFail(t *testing.T) {
	t.Parallel()

	tableName := "test-table"

	// Create a mock instance of the DynamoDB interface
	mockDynamoDB := ddb.NewMockDynamodbAPI(t)

	// Create a new store using the mock DynamoDB instance
	store := ddb.NewStore(tableName, mockDynamoDB)

	// Set up the expected inputs and outputs
	ctx := context.Background()
	bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")
	keyItem := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{Value: bookID.String()},
	}
	expectedGetItemInput := &dynamodb.GetItemInput{
		Key:       keyItem,
		TableName: aws.String(tableName),
	}

	// Expect a call to GetItem with the expected input and return the mock output
	mockDynamoDB.EXPECT().GetItem(ctx, expectedGetItemInput).Return(&dynamodb.GetItemOutput{}, assert.AnError).Once()

	// Call the FindOne method of the store
	_, err := store.FindOne(ctx, bookID)

	// Assert the expected output
	require.Error(t, err)

	// Assert that the mockDynamoDB's expectations were met
	mockDynamoDB.AssertExpectations(t)
}
