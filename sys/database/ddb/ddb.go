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

// DefaultTableScanLimit is the default limit for the Scan operation.
//
// NOTE: This is a temporary solution to avoid scanning the entire table.
const DefaultTableScanLimit = 25

// DynamodbAPI is the interface used to interact with AWS DynamoDB.
//
//go:generate mockery --name DynamoDB
type DynamodbAPI interface {
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
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

// FindAll returns all books from the DynamoDB database.
func (s *Store) FindAll(ctx context.Context) ([]domain.Book, error) {
	response, err := s.api.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(s.tableName),
		Limit:     aws.Int32(DefaultTableScanLimit),
	})

	if err != nil {
		return []domain.Book{}, fmt.Errorf("findall: scan: %w", err)
	}

	items := make([]DynamodbBook, 0, len(response.Items))

	if err = attributevalue.UnmarshalListOfMaps(response.Items, &items); err != nil {
		return []domain.Book{}, fmt.Errorf("findall: unmarshallistofmaps: %w", err)
	}

	books := ToDomainBooks(items)

	return books, nil
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

	if len(response.Item) == 0 {
		return domain.Book{}, fmt.Errorf("findbyid: getitem[%v]: %w", bookID, domain.ErrNotFound)
	}

	if err = attributevalue.UnmarshalMap(response.Item, &item); err != nil {
		return domain.Book{}, fmt.Errorf("findbyid: unmarshalmap: %w", err)
	}

	book := ToDomainBook(item)

	return book, nil
}

// Delete removes a book from the DynamoDB database by using bookID as primary key.
func (s *Store) Delete(ctx context.Context, bookID uuid.UUID) error {
	_, err := s.api.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: bookID.String()},
		},
		ConditionExpression: aws.String("attribute_exists(id)"),
	})

	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
