package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/rotiroti/alessandrina/sys/database/memory"
	"github.com/rotiroti/alessandrina/web"
)

func main() {
	app, err := run()
	if err != nil {
		log.Fatal(err)
	}

	lambda.Start(app.CreateBook)
}

func run() (*web.APIGatewayV2Handler, error) {
	// Instantiate a new memory store
	store := memory.NewStore()

	// Instantiate a new domain service
	service := domain.NewService(store)

	// Instantiate a new handler
	handler := web.NewAPIGatewayV2Handler(service)

	return handler, nil
}
