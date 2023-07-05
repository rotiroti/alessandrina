package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/rotiroti/alessandrina/sys/database/ddb"
	"github.com/rotiroti/alessandrina/sys/database/memory"
	"github.com/rotiroti/alessandrina/web"
)

var (
	// ErrMissingTableName is returned when the TABLE_NAME environment variable is not set.
	ErrMissingTableName = errors.New("missing TABLE_NAME environment variable")
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("startup: %v\n", err)
	}
}

func run(ctx context.Context) error {
	type Storage int

	const (
		// DynamoDB will store data in AWS DynamoDB.
		DynamoDB Storage = iota

		// Memory will store data in memory.
		Memory
	)

	// Set the default storage to DynamoDB.
	defaultStorage := DynamoDB

	// Check if the storage is set to memory.
	if os.Getenv("STORAGE_MEMORY") == "true" {
		defaultStorage = Memory
	}

	var store domain.Storer

	switch defaultStorage {
	case Memory:
		// Instantiate a new memory store
		store = memory.NewStore()
	case DynamoDB:
		tableName, ok := os.LookupEnv("TABLE_NAME")
		if !ok || tableName == "" {
			return fmt.Errorf("run: %w", ErrMissingTableName)
		}

		// Define a custom endpoint resolver to use a local DynamoDB instance.
		// This is useful for local development, for example with the SAM CLI and LocalStack.
		customResolver := aws.EndpointResolverWithOptionsFunc(
			func(_, region string, options ...interface{}) (aws.Endpoint, error) {
				endpoint := os.Getenv("AWS_ENDPOINT_DEBUG")
				if endpoint == "" {
					return aws.Endpoint{}, &aws.EndpointNotFoundError{}
				}

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

		// TODO: This is a temporary solution to check if we can connect to DynamoDB.
		//	 We should declare a Config struct in the ddb package, and use it
		//	 instead of instantiating the DynamoDB client here.
		conf, err := config.LoadDefaultConfig(ctx, options...)
		if err != nil {
			return err
		}

		// Instantiate a new DynamoDB store
		store = ddb.NewStore(tableName, dynamodb.NewFromConfig(conf))
	}

	// Instantiate a new domain service
	service := domain.NewService(store)

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(service)

	// Start the lambda handler listening for CreateBook events.
	lambda.Start(handler.CreateBook)

	return nil
}
