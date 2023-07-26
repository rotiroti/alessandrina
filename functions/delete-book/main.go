package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/rotiroti/alessandrina/sys/database/ddb"
	"github.com/rotiroti/alessandrina/web"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatalf("startup: %v\n", err)
	}
}

func run(ctx context.Context) error {
	var (
		store *ddb.Store
		err   error
	)

	tableName := os.Getenv("TABLE_NAME")
	debugMode := os.Getenv("DEBUG_MODE")

	switch debugMode {
	case "true":
		store, err = ddb.NewDebugStore(ctx, tableName)
	default:
		store, err = ddb.NewStore(ctx, tableName)
	}

	if err != nil {
		return err
	}

	bookCore := domain.NewBookCore(store)
	handler := web.NewAPIGatewayV2Handler(bookCore)
	lambda.Start(handler.DeleteBook)

	return nil
}
