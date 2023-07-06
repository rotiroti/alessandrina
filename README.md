# Alessandrina

This project aims to build a Go-based serverless application using AWS SAM. It provides an API with endpoints for interacting with a book database, allowing users to search, create, and delete books. The project also includes a robust CI/CD pipeline for automated build, test, and deployment on AWS using GitHub Actions.

## Requirements

To set up and run this serverless application locally or in a cloud environment, ensure you have the following prerequisites:

- [AWS Account](https://aws.amazon.com/account)
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html)
- [Go Programming Language](https://go.dev)
- [GNU Make](https://www.gnu.org/software/make)
- [Docker](https://www.docker.com)
- [Artillery](https://artillery.io)
- [Localstack](https://localstack.cloud) (required only for running AWS DynamoDB locally)

Once you have these prerequisites, you can set up and run the serverless application locally or deploy it to your preferred cloud environment.

## Project Structure

The serverless application is structured using the *hexagonal architecture* known as *ports and adapters*. This architectural pattern provides a way to separate the core business logic (`domain`) of the application from specific technical implementations, like the infrastructure (`database`) and the handling of client requests (`web`). Encapsulating the domain logic within the hexagon makes it easier to maintain and modify the application without affecting other components.

```shell
├── domain
│  ├── book.go
│  ├── book_test.go
│  ├── mock_storer_test.go
│  ├── mock_uuid_generator_test.go
│  └── model.go
├── events
│  └── create.json
├── functions
│  └── create-book
│     └── main.go
├── go.mod
├── go.sum
├── Makefile
├── README.md
├── samconfig.toml
├── scripts
│  └── create-table.sh
├── sys
│  └── database
│     ├── ddb
│     │  ├── ddb.go
│     │  ├── ddb_test.go
│     │  ├── mock_dynamodb_api_test.go
│     │  └── model.go
│     └── memory
│        ├── memory.go
│        └── memory_test.go
├── template.yaml
├── tests
│  ├── integration
│  │  └── main_test.go
│  └── performance
│     ├── books50.csv
│     ├── load.yml
│     └── spike.yml
└── web
   ├── apigateway.go
   ├── apigateway_test.go
   └── model.go
```

### `/functions`

This project's entire AWS Lambda functions inside the `/functions` folder. The folders under `/functions` are consistently named for each lambda that SAM will build. Each folder has a matching source code file that contains the `main` package. None of the packages inside the folder `functions` can import each other.

### `/tests`

The `/tests` folder houses a collection of `integration` and `performance` tests specifically designed to evaluate the functionality and performance of the serverless application. These tests are developed to simulate real-world scenarios and interactions with the application, ensuring it behaves as expected and performs optimally under different conditions.

### `/events`

This folder contains multiple event bodies that can be passed to SAM when invoking the AWS serverless functions locally.

```shell
sam local invoke -e events/<EVENT_NAME>.json
```

### `/scripts`

This folder contains shell scripts to perform migrations when running DynamoDB on Localstack.

## Environment Variables

The serverless application can be configured via some environment variables.

```shell
# Use the in-memory database (sys/database/memory package)
STORAGE_MEMORY: false

# Set the table name (required when using DynamoDB storage)
TABLE_NAME: BooksTable-local

# Set a custom endpoint to be used for a service (i.e. Localstack)
AWS_ENDPOINT_DEBUG: "http://localstack_main:4566"

# Enable debug logging for AWS SDK Go clients to view HTTP requests and response bodies.
AWS_CLIENT_DEBUG: false
```

## Makefile Commands

```shell
# Perform unit tests.
make unit-tests

# Run unit tests and create an HTML code coverage report.
make coverage

# Format Go source files.
make format

# Run Go linters aggregator (golangci-lint) using Docker.
make lint

# Perform integration tests. (see Integration Tests setup section)
make integration-tests

# Generate Go mocks
make mocks

# Remove Go mocks
make remove-mocks

# Cleanup artifacts, coverage and report files.
make clean
```

## Integration Tests

### In-memory database

This setup assume you have alredy installed all the previous requirements.

```shell
# 1. Build the serverless application.
sam build --debug

# 2. Start a local HTTP API server using a mock (in-memory) database.
STORAGE_MEMORY=true sam local start-api --debug

# 3. Open another shell and execute the command in the project's root directory.
make integration-tests
```

### AWS DynamoDB (Localstack)

These stesp assume that you have already installed all the requirements mentioned in the "Requirements" section.

```shell
# 1. Create a Docker network
docker create network alessandrina

# 2. Start the Localstack server.
DOCKER_FLAGS="--network alessandrina -d" localstack start

# 3. Create a new DynamoDB table on Localstack
sh ./scripts/create-table.sh BooksTable-local

# 4. Build the serverless application.
sam build --debug

# 5. Start a local HTTP API server using the "BooksTable-local" DynamoDB table.
TABLE_NAME=BooksTable-local AWS_ENDPOINT_DEBUG="http://localstack_main:4566" \
    sam local start-api \
        --docker-network alessandrina \
        --warm-containers LAZY \
        --debug

# 6. Open another shell and execute the command in the project's root directory.
make integration-tests
```
