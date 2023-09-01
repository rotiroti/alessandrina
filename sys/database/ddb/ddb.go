package ddb

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
)

// ErrMissingTableName is returned when the DB_TABLE environment variable is not set.
var ErrMissingTableName = errors.New("missing DB_TABLE environment variable")

// DefaultTableScanLimit is the default limit for the Scan operation.
//
// NOTE: This is a temporary solution to avoid scanning the entire table.
const DefaultTableScanLimit = 25

// Option is a function that configures a Store.
type Option func(*Store) error

// WithClient returns a Store Option that sets the DynamoDB client.
func WithClient(client DynamoDBClient) Option {
	return func(s *Store) error {
		s.client = client

		return nil
	}
}

// WithClientLog returns a Store Option that sets the DynamoDB client with logging enabled.
func WithClientLog() Option {
	return func(s *Store) error {
		logMode := aws.LogRequestWithBody | aws.LogResponseWithBody
		conf, err := config.LoadDefaultConfig(context.Background(), config.WithClientLogMode(logMode))

		if err != nil {
			return fmt.Errorf("ddb.withclientlog loaddefaultconfig: %w", err)
		}

		s.client = dynamodb.NewFromConfig(conf)

		return nil
	}
}

// WithLocalStack returns a Store Option that sets the DynamoDB client with LocalStack.
func WithLocalStack() Option {
	return func(s *Store) error {
		options := []func(*config.LoadOptions) error{}
		resolver := aws.EndpointResolverWithOptionsFunc(
			func(_, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           "http://localstack_main:4566",
					SigningRegion: region,
				}, nil
			})
		logMode := aws.LogRequestWithBody | aws.LogResponseWithBody
		options = append(options,
			config.WithEndpointResolverWithOptions(resolver),
			config.WithClientLogMode(logMode),
		)

		conf, err := config.LoadDefaultConfig(context.Background(), options...)
		if err != nil {
			return fmt.Errorf("ddb.withlocalstack loaddefaultconfig: %w", err)
		}

		s.client = dynamodb.NewFromConfig(conf)

		return nil
	}
}

// DynamoDBClient is the interface used to interact with AWS DynamoDB.
type DynamoDBClient interface {
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

// Store is a DynamoDB implementation of the Storer interface.
type Store struct {
	client DynamoDBClient
	table  string
}

// Ensure Store implements the Storer interface.
var _ domain.Storer = (*Store)(nil)

// NewStore returns a new DynamoDB Store.
func NewStore(ctx context.Context, table string, opts ...Option) (*Store, error) {
	if table == "" {
		return nil, fmt.Errorf("ddb.newstore: %w", ErrMissingTableName)
	}

	store := &Store{table: table}

	for _, opt := range opts {
		err := opt(store)
		if err != nil {
			return nil, fmt.Errorf("ddb.newstore: %w", err)
		}
	}

	if store.client == nil {
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			return nil, fmt.Errorf("ddb.newstore loaddefaultconfig: %w", err)
		}

		store.client = dynamodb.NewFromConfig(cfg)
	}

	return store, nil
}

// Save adds a new book into the DynamoDB database.
func (s *Store) Save(ctx context.Context, book domain.Book) error {
	item, err := attributevalue.MarshalMap(ToDynamodbBook(book))
	if err != nil {
		return fmt.Errorf("ddb.save marshalmap: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(s.table),
	})

	if err != nil {
		return fmt.Errorf("ddb.save putitem: %w", err)
	}

	return nil
}

// FindAll returns all books from the DynamoDB database.
func (s *Store) FindAll(ctx context.Context) ([]domain.Book, error) {
	response, err := s.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(s.table),
		Limit:     aws.Int32(DefaultTableScanLimit),
	})

	if err != nil {
		return []domain.Book{}, fmt.Errorf("ddb.findall scan: %w", err)
	}

	items := make([]DynamodbBook, 0, len(response.Items))

	if err = attributevalue.UnmarshalListOfMaps(response.Items, &items); err != nil {
		return []domain.Book{}, fmt.Errorf("ddb.findall unmarshallistofmaps: %w", err)
	}

	books := ToDomainBooks(items)

	return books, nil
}

// FindOne returns a book from the DynamoDB database by using bookID as primary key.
func (s *Store) FindOne(ctx context.Context, bookID uuid.UUID) (domain.Book, error) {
	item := DynamodbBook{ID: bookID.String()}
	response, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: item.ID},
		},
	})

	if err != nil {
		return domain.Book{}, fmt.Errorf("ddb.findone getitem: %w", err)
	}

	if len(response.Item) == 0 {
		return domain.Book{}, fmt.Errorf("ddb.findone getitem: %w", domain.ErrNotFound)
	}

	if err = attributevalue.UnmarshalMap(response.Item, &item); err != nil {
		return domain.Book{}, fmt.Errorf("ddb.findone unmarshalmap: %w", err)
	}

	book := ToDomainBook(item)

	return book, nil
}

// Delete removes a book from the DynamoDB database by using bookID as primary key.
func (s *Store) Delete(ctx context.Context, bookID uuid.UUID) error {
	_, err := s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: bookID.String()},
		},
		ConditionExpression: aws.String("attribute_exists(id)"),
	})

	if err != nil {
		return fmt.Errorf("ddb.delete deleteitem: %w", err)
	}

	return nil
}
