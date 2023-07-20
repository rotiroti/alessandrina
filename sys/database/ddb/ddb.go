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

// ErrMissingTableName is returned when the TABLE_NAME environment variable is not set.
var ErrMissingTableName = errors.New("missing TABLE_NAME environment variable")

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

// Config is the configuration used to create a new DynamoDB client.
type Config struct {
	ClientLog string
	Endpoint  string
	TableName string
}

func parse(ctx context.Context, conf Config) ([]func(*config.LoadOptions) error, error) {
	options := []func(*config.LoadOptions) error{}

	if conf.TableName == "" {
		return options, fmt.Errorf("parse: %w", ErrMissingTableName)
	}

	// Define a custom endpoint resolver to use a local DynamoDB instance.
	// This is useful for local development, for example with the SAM CLI and LocalStack.
	resolver := aws.EndpointResolverWithOptionsFunc(
		func(_, region string, options ...interface{}) (aws.Endpoint, error) {
			endpoint := conf.Endpoint
			if endpoint != "" {
				return aws.Endpoint{
					PartitionID:   "aws",
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			}

			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

	// Enable debug logging to see the HTTP requests and responses bodies.
	var logMode aws.ClientLogMode

	if conf.ClientLog == "true" {
		logMode |= aws.LogRequestWithBody | aws.LogResponseWithBody
	}

	options = append(options,
		config.WithEndpointResolverWithOptions(resolver),
		config.WithClientLogMode(logMode),
	)

	return options, nil
}

// Store is a DynamoDB implementation of the Storer interface.
type Store struct {
	Client DynamodbAPI
	config Config
}

// Ensure Store implements the Storer interface.
var _ domain.Storer = (*Store)(nil)

// NewStore returns a new instance of Store.
func NewStore(ctx context.Context, conf Config) (*Store, error) {
	options, err := parse(ctx, conf)
	if err != nil {
		return &Store{}, fmt.Errorf("newstore: %w", err)
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		return &Store{}, fmt.Errorf("newstore: %w", err)
	}

	return &Store{
		Client: dynamodb.NewFromConfig(awsConfig),
		config: conf,
	}, nil
}

// Save adds a new book into the DynamoDB database.
func (s *Store) Save(ctx context.Context, book domain.Book) error {
	item, err := attributevalue.MarshalMap(ToDynamodbBook(book))
	if err != nil {
		return fmt.Errorf("marshalmap: %w", err)
	}

	_, err = s.Client.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(s.config.TableName),
	})

	if err != nil {
		return fmt.Errorf("putitem: %w", err)
	}

	return nil
}

// FindAll returns all books from the DynamoDB database.
func (s *Store) FindAll(ctx context.Context) ([]domain.Book, error) {
	response, err := s.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(s.config.TableName),
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
	response, err := s.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.config.TableName),
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
	_, err := s.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(s.config.TableName),
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
