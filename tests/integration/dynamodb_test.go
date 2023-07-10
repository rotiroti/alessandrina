package integration

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/rotiroti/alessandrina/sys/database/ddb"
	"github.com/stretchr/testify/require"
)

func skipLocalstack(t *testing.T) {
	t.Helper()
	if os.Getenv("LOCALSTACK") == "" {
		t.Skip("skipping integration tests, set environment variable LOCALSTACK")
	}
}

func setupLocalstack(t *testing.T) (tableName string, client *dynamodb.Client) {
	t.Helper()

	ctx := context.Background()

	tableName, ok := os.LookupEnv("TABLE_NAME")
	if !ok || tableName == "" {
		t.Fatalf("TABLE_NAME environment variable not set")
	}

	customResolver := aws.EndpointResolverWithOptionsFunc(
		func(_, region string, options ...interface{}) (aws.Endpoint, error) {
			endpoint := "http://localstack_main:4566"

			return aws.Endpoint{
				PartitionID:       "aws",
				URL:               endpoint,
				SigningRegion:     region,
				HostnameImmutable: true,
			}, nil
		})

	// Enable debug logging to see the HTTP requests and responses bodies.
	var logMode aws.ClientLogMode

	if os.Getenv("AWS_CLIENT_DEBUG") == "true" {
		logMode |= aws.LogRequestWithBody | aws.LogResponseWithBody
	}

	options := []func(*config.LoadOptions) error{
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithClientLogMode(logMode),
	}

	conf, err := config.LoadDefaultConfig(ctx, options...)
	if err != nil {
		t.Fatalf("failed to load SDK configuration, %v", err)
	}

	client = dynamodb.NewFromConfig(conf)

	return tableName, client
}

func TestDynamodbFlow(t *testing.T) {
	// t.Parallel()

	os.Setenv("LOCALSTACK", "true")
	os.Setenv("TABLE_NAME", "BooksTable-local")

	// Skip the DynamoDB test if the LOCALSTACK environment variable is not set
	skipLocalstack(t)

	// Setup test environment
	ctx := context.Background()
	tableName, client := setupLocalstack(t)

	t.Run("Save", func(t *testing.T) {
		// t.Parallel()

		// Create a new book
		bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1c812")
		book := domain.Book{
			ID:        bookID,
			Title:     "The Lord of the Rings",
			Authors:   "J.R.R. Tolkien",
			Pages:     1178,
			Publisher: "George Allen & Unwin",
			ISBN:      "978-0-261-10235-4",
		}

		// Create a new store using the Localstack DynamoDB instance
		//
		// NOTE: the table name should coincide with the one created
		//       by using the `scripts/create-table` script
		store := ddb.NewStore(tableName, client)

		// Call the Save method of the store
		err := store.Save(ctx, book)

		// Assert the expected output
		require.NoError(t, err)
	})

	t.Run("SaveTableNotExist", func(t *testing.T) {
		// t.Parallel()

		// Create a new store using the Localstack DynamoDB instance
		store := ddb.NewStore("not-exist", client)

		// Call the Save method of the store
		err := store.Save(ctx, domain.Book{})

		// Assert the expected output
		require.Error(t, err)
	})

	t.Run("FindOne", func(t *testing.T) {
		// t.Parallel()

		// Create a new book
		bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1cabd")
		expectedBook := domain.Book{
			ID:        bookID,
			Title:     "The Lord of the Rings",
			Authors:   "J.R.R. Tolkien",
			Pages:     1178,
			Publisher: "George Allen & Unwin",
			ISBN:      "978-0-261-10235-4",
		}

		// Create a new store using the Localstack DynamoDB instance
		store := ddb.NewStore(tableName, client)

		// Call the Save method of the store
		err := store.Save(ctx, expectedBook)

		// Assert the expected output
		require.NoError(t, err)

		// Call the FindOne method of the store
		book, err := store.FindOne(ctx, bookID)

		// Assert the expected output
		require.NoError(t, err)

		// Assert the expected output
		require.Equal(t, expectedBook, book)
	})

	t.Run("FindOneBookNotFound", func(t *testing.T) {
		// t.Parallel()

		// Create a new book
		bookID := uuid.MustParse("01234567-0123-0123-0123-0123456789ab")

		// Create a new store using the Localstack DynamoDB instance
		store := ddb.NewStore(tableName, client)

		// Call the FindOne method of the store
		_, err := store.FindOne(ctx, bookID)

		// Assert the expected output
		require.Error(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		// t.Parallel()

		// Create a new book
		bookID := uuid.MustParse("ad8b59c2-5fe6-4267-b0cf-6d2f9eb1cabd")
		expectedBook := domain.Book{
			ID:        bookID,
			Title:     "The Lord of the Rings",
			Authors:   "J.R.R. Tolkien",
			Pages:     1178,
			Publisher: "George Allen & Unwin",
			ISBN:      "978-0-261-10235-4",
		}

		// Create a new store using the Localstack DynamoDB instance
		store := ddb.NewStore(tableName, client)

		// Call the Save method of the store
		err := store.Save(ctx, expectedBook)

		// Assert the expected output
		require.NoError(t, err)

		// Call the Delete method of the store
		err = store.Delete(ctx, bookID)

		// Assert the expected output
		require.NoError(t, err)
	})

	t.Run("DeleteBookNotFound", func(t *testing.T) {
		// t.Parallel()

		// Create a new book
		bookID := uuid.MustParse("01234567-0123-0123-0123-0123456789ab")

		// Create a new store using the Localstack DynamoDB instance
		store := ddb.NewStore(tableName, client)

		// Call the Delete method of the store
		err := store.Delete(ctx, bookID)

		// Assert the expected output
		require.Error(t, err)
	})
}
