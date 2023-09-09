API_URL = http://localhost:3000
CUR_DIR = $(shell echo "${PWD}")

.PHONY: unit-tests integration-tests coverage format clean mocks remove-mocks lint

build-HelloFunction:
	@echo "Building HelloFunction"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o hello github.com/rotiroti/alessandrina/functions/hello/
	mv hello $(ARTIFACTS_DIR)
	@echo "Built HelloFunction successfully"

build-GetBooksFunction:
	@echo "Building GetBooksFunction"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o get-books github.com/rotiroti/alessandrina/functions/get-books/
	mv get-books $(ARTIFACTS_DIR)
	@echo "Built GetBooksFunction successfully"

build-GetBookFunction:
	@echo "Building GetBookFunction"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o get-book github.com/rotiroti/alessandrina/functions/get-book/
	mv get-book $(ARTIFACTS_DIR)
	@echo "Built GetBookFunction successfully"

build-CreateBookFunction:
	@echo "Building CreateBookFunction"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o create-book github.com/rotiroti/alessandrina/functions/create-book/
	mv create-book $(ARTIFACTS_DIR)
	@echo "Built CreateBookFunction successfully"

build-DeleteBookFunction:
	@echo "Building DeleteBookFunction"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o delete-book github.com/rotiroti/alessandrina/functions/delete-book/
	mv delete-book $(ARTIFACTS_DIR)
	@echo "Built DeleteBookFunction successfully"

build-k6:
	@echo "Building k6 tool"
	go install go.k6.io/xk6/cmd/xk6@latest
	xk6 build --with github.com/grafana/xk6-dashboard@latest
	@echo "Built k6 successfully"

format:
	@echo "Format code"
	go fmt ./...

lint:
	@echo "Run linting"
	@docker run -t --rm -v ${CUR_DIR}:/app -w /app golangci/golangci-lint:v1.53 golangci-lint run --fix -v

unit-tests:
	@echo "Run unit tests"
	go test -v -race -coverprofile=coverage.out -covermode=atomic $$(go list ./... | grep -v /functions/ | grep -v /tests/)

coverage: unit-tests
	@echo "Run unit tests and create HTML coverage report"
	go tool cover -html=coverage.out -o coverage.html

integration-tests:
	@echo "Run integration tests"
	API_URL=${API_URL} INTEGRATION=1 go test -count=1 -v -race ./tests/...

mocks:
	@echo "Generate mocks"
	docker run --rm -v ${CUR_DIR}:/app -w /app vektra/mockery --all

remove-mocks:
	@echo "Clean up mocks"
	find . -name 'mock_*.go' | xargs rm
	rm -fr mocks/

clean:
	@echo "Clean up"
	rm -f coverage.*
	rm -f report.*
	rm -fr .aws-sam
	rm -f ./k6
