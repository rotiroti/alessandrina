package ddb

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
)

// DynamodbAPI is the interface used to interact with AWS DynamoDB.
//
//go:generate mockery --name DynamoDB
type DynamodbAPI interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

// Store is a DynamoDB implementation of the Storer interface.
type Store struct {
	api       DynamodbAPI
	tableName string
}

// Ensure Store implements the Storer interface.
var _ domain.Storer = (*Store)(nil)

// NewStore returns a new instance of Store.
func NewStore(tableName string, api DynamodbAPI) *Store {
	return &Store{
		tableName: tableName,
		api:       api,
	}
}

// Save adds a new book into the DynamoDB database.
func (s *Store) Save(ctx context.Context, book domain.Book) error {
	item, err := attributevalue.MarshalMap(ToDynamodbBook(book))
	if err != nil {
		return fmt.Errorf("marshalmap: %w", err)
	}

	_, err = s.api.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(s.tableName),
	})

	if err != nil {
		return fmt.Errorf("putitem: %w", err)
	}

	return nil
}

// FindOne returns a book from the DynamoDB database by using bookID as primary key.
func (s *Store) FindOne(ctx context.Context, bookID uuid.UUID) (domain.Book, error) {
	item := DynamodbBook{ID: bookID.String()}
	response, err := s.api.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: item.ID},
		},
	})

	if err != nil {
		return domain.Book{}, fmt.Errorf("findbyid: getitem[%v]: %w", bookID, err)
	}

	if err = attributevalue.UnmarshalMap(response.Item, &item); err != nil {
		return domain.Book{}, fmt.Errorf("findbyid: unmarshalmap: %w", err)
	}

	book := ToDomainBook(item)

	return book, nil
}
