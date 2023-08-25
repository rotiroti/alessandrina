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

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func run(ctx context.Context) error {
	dbTable := getEnv("DB_TABLE", "")
	dbConn := getEnv("DB_CONNECTION", "aws")
	dbLog := getEnv("DB_LOG", "false")

	var (
		store *ddb.Store
		err   error
	)

	switch dbConn {
	case "localstack":
		store, err = ddb.NewStore(ctx, dbTable, ddb.WithLocalStack())
	default:
		if dbLog == "true" {
			store, err = ddb.NewStore(ctx, dbTable, ddb.WithClientLog())
		} else {
			store, err = ddb.NewStore(ctx, dbTable)
		}
	}

	if err != nil {
		return err
	}

	bookCore := domain.NewBookCore(store)
	handler := web.NewAPIGatewayV2Handler(bookCore)

	lambda.Start(handler.DeleteBook)

	return nil
}
