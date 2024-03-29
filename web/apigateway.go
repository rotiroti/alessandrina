package web

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/rotiroti/alessandrina/domain"
	"github.com/rotiroti/alessandrina/sys/validation"
)

// APIGatewayV2Handler is the handler for the API Gateway v2.
type APIGatewayV2Handler struct {
	book      *domain.BookCore
	validator validation.Validator
}

// NewAPIGatewayV2Handler returns a new APIGatewayV2Handler.
func NewAPIGatewayV2Handler(book *domain.BookCore) *APIGatewayV2Handler {
	return &APIGatewayV2Handler{
		book:      book,
		validator: validation.New(),
	}
}

// CreateBook handles requests for creating a book.
func (h *APIGatewayV2Handler) CreateBook(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var appNewBook AppNewBook

	if err := json.Unmarshal([]byte(req.Body), &appNewBook); err != nil {
		return errorResponse(http.StatusBadRequest, err.Error()), nil
	}

	if err := h.validator.Check(appNewBook); err != nil {
		return errorResponse(http.StatusBadRequest, err.Error()), nil
	}

	domainNewBook := ToDomainNewBook(appNewBook)
	ret, err := h.book.Save(ctx, domainNewBook)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return jsonResponse(http.StatusCreated, ToAppBook(ret)), nil
}

// GetBooks handles requests for getting all books.
func (h *APIGatewayV2Handler) GetBooks(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ret, err := h.book.FindAll(ctx)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return jsonResponse(http.StatusOK, ToAppListBooks(ret)), nil
}

// GetBook handles requests for getting a book by a given ID (UUID).
func (h *APIGatewayV2Handler) GetBook(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	id, err := uuid.Parse(req.PathParameters["id"])
	if err != nil {
		return errorResponse(http.StatusBadRequest, err.Error()), nil
	}

	ret, err := h.book.FindOne(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return errorResponse(http.StatusNotFound, err.Error()), nil
		}

		return errorResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return jsonResponse(http.StatusOK, ToAppBook(ret)), nil
}

// DeleteBook handles requests for deleting a book by a given ID (UUID).
func (h *APIGatewayV2Handler) DeleteBook(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	id, err := uuid.Parse(req.PathParameters["id"])
	if err != nil {
		return errorResponse(http.StatusBadRequest, err.Error()), nil
	}

	if err := h.book.Delete(ctx, id); err != nil {
		return errorResponse(http.StatusInternalServerError, err.Error()), nil
	}

	return jsonResponse(http.StatusNoContent, nil), nil
}

func jsonResponse(code int, obj any) events.APIGatewayV2HTTPResponse {
	body, err := json.Marshal(obj)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, err.Error())
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(body),
		IsBase64Encoded: false,
	}
}

func errorResponse(code int, err string) events.APIGatewayV2HTTPResponse {
	type data struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	errorMessage := map[string]data{
		"error": {
			Code:    code,
			Message: err,
		},
	}

	// NOTE: ignoring error as if Marshal fails even here, we have bigger problems.
	body, _ := json.Marshal(errorMessage)

	return events.APIGatewayV2HTTPResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            string(body),
		IsBase64Encoded: false,
	}
}
