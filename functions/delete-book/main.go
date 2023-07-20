package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/rotiroti/alessandrina/sys/database/ddb"
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
	DB := ddb.Config{
		TableName: os.Getenv("TABLE_NAME"),
		Endpoint:  os.Getenv("AWS_ENDPOINT_DEBUG"),
		ClientLog: os.Getenv("AWS_CLIENT_DEBUG"),
	}

	store, err := ddb.NewStore(ctx, DB)
	if err != nil {
		return err
	}

	bookCore := domain.NewBookCore(store)
	handler := web.NewAPIGatewayV2Handler(bookCore)
	lambda.Start(handler.DeleteBook)

	return nil
}
