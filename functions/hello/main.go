package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type book struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Autors    string `json:"autors"`
	ISBN      string `json:"isbn"`
	Publisher string `json:"publisher"`
	Page      int    `json:"page"`
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	b := &book{
		ID:        "1",
		Title:     "The Go Programming Language",
		Autors:    "Alan Donovan, Brian W. Kernighan",
		ISBN:      "9780134190440",
		Publisher: "Addison-Wesley Professional Computing Series",
		Page:      400,
	}

	body, err := json.Marshal(b)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error: %s", err),
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, err
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
