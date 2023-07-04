AWS_TESTING_API_URL = http://localhost:3000
CUR_DIR = $(shell echo "${PWD}")

.PHONY: unit-tests integration-tests coverage format clean mocks remove-mocks lint

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
	AWS_TESTING_API_URL=${AWS_TESTING_API_URL} INTEGRATION=1 go test -count=1 -v -race ./tests/...

build-CreateBookFunction:
	@echo "Building CreateBookFunction"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o create-book github.com/rotiroti/alessandrina/functions/create-book/
	mv create-book $(ARTIFACTS_DIR)
	@echo "Built CreateBookFunction successfully"

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
	rm -fr .aws-sam
