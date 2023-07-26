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
	store, err := ddb.NewStore(ctx, os.Getenv("TABLE_NAME"))
	if err != nil {
		return err
	}
	bookCore := domain.NewBookCore(store)
	handler := web.NewAPIGatewayV2Handler(bookCore)
	lambda.Start(handler.DeleteBook)

	return nil
}
