# Alessandrina

[![codecov](https://codecov.io/gh/rotiroti/alessandrina/branch/main/graph/badge.svg?token=eWAHfGU54Y)](https://codecov.io/gh/rotiroti/alessandrina)
![CI/CD](https://github.com/rotiroti/alessandrina/actions/workflows/pipeline.yaml/badge.svg)

This project aims to build a Go-based serverless application using AWS SAM. It provides an API with endpoints for interacting with a book database, allowing users to search, create, and delete books. The project also includes a robust CI/CD pipeline for automated build, test, and deployment on AWS using GitHub Actions.

## Requirements

To set up and run this serverless application locally or in a cloud environment, ensure you have the following prerequisites:

- [AWS Account](https://aws.amazon.com/account)
- [AWS SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html)
- [Go Programming Language](https://go.dev)
- [GNU Make](https://www.gnu.org/software/make)
- [Docker](https://www.docker.com)
- [k6](https://k6.io/)
- [Localstack](https://localstack.cloud) (required only for running AWS DynamoDB locally)

Once you have these prerequisites, you can set up and run the serverless application locally or deploy it to your preferred cloud environment.

## Serverless Architecture

<p align="center">
  <img src="assets/architecture.png" alt="Alessandrina Architecture"/>
</p>

## Project Structure

The serverless application is structured using the *hexagonal architecture* known as *ports and adapters*. This architectural pattern provides a way to separate the core business logic (`domain`) of the application from specific technical implementations, like the infrastructure (`database`) and the handling of client requests (`web`). Encapsulating the domain logic within the hexagon makes it easier to maintain and modify the application without affecting other components.

```shell
├── assets
├── domain
├── events
├── functions
│  ├── create-book
│  ├── delete-book
│  ├── get-book
│  └── get-books
├── go.mod
├── go.sum
├── locals.json
├── Makefile
├── README.md
├── samconfig.toml
├── scripts
│  ├── create-table.sh
│  └── delete-table.sh
├── sys
│  └── database
│     ├── ddb
│     └── memory
├── template.yaml
├── tests
│  ├── integration
│  └── performance
└── web
```

### `/functions`

This project's entire AWS Lambda functions inside the `/functions` folder. The folders under `/functions` are consistently named for each lambda that SAM will build. Each folder has a matching source code file that contains the `main` package. None of the packages inside the folder `functions` can import each other.

### `/tests`

The `/tests` folder houses a collection of `integration` and `performance` tests specifically designed to evaluate the functionality and performance of the serverless application. These tests are developed to simulate real-world scenarios and interactions with the application, ensuring it behaves as expected and performs optimally under different conditions.

### `/events`

This folder contains multiple event bodies that can be passed to SAM when invoking the AWS serverless functions locally.

```shell
sam local invoke CreateBookFunction \
   -e events/create-book.json \
   --docker-network alessandrina \
   --env-vars locals.json
```

### `/scripts`

This folder contains shell scripts to perform migrations when running DynamoDB on Localstack.

## Environment Variables for SAM

The serverless application can be configured via some environment variables.

```shell
# Set the table name (mandatory)
DB_TABLE=BooksTable-local

# Set the DynamoDB client connection (possible values: aws|localstack, default: aws)
DB_CONNECTION=localstack

# Enable AWS Client Logs for the DynamoDB service (default: false)
#
# When running with DB_CONNECTION=localstack, client logs are enabled as default
DB_LOG=true
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

The integration tests assume that you have already installed all the requirements mentioned in the "Requirements" section.

### Running on SAM Local API + Localstack DynamoDB

```shell
# 1. Create a Docker network
docker create network alessandrina

# 2. Start the Localstack server.
DOCKER_FLAGS="--network alessandrina -d" localstack start

# 3. Create a new DynamoDB table on Localstack
sh ./scripts/create-table.sh BooksTable-local

# 4. Build the serverless application.
sam build --parallel

# 5. Start a local HTTP API server using the "BooksTable-local" DynamoDB table.
sam local start-api --docker-network alessandrina --warm-containers LAZY --env-vars locals.json

# 6. Open another shell and execute the command in the project's root directory.
make integration-tests
```

### Running on AWS (feature, dev, prod branches)

```shell
# 1. Build the serverless application
sam build --parallel

# 2. Run integration tests
make integration-tests API_URL=<STACK_WEBPOINT_URL>
```

## Performance Tests

The performance tests assume that you have already installed all the requirements mentioned in the "Requirements" section.

### Environment Variables for k6

```shell
# Set the API URL (mandatory)
API_URL=<STACK_WEBPOINT_URL>

# Set the workload (possible values: 0|1|2, default: 0)
WORKLOAD=0

# Set the book operation (mandatory, possible values: list|create|flow)
BOOK_OP=list

# Set the test name (default: main.js)
TEST_NAME=smoke-create
```

### Run a test locally

```shell
API_URL=<STACK_WEBPOINT_URL> BOOK_OP=list k6 run ./tests/performance/main.js
```

#### Authenticate with Grafana Cloud

```shell
k6 login cloud --token <PERSONAL_API_TOKEN>
```

### Run a test in the cloud

```shell
./k6 cloud -e K6_CLOUD_PROJECT_ID=<PROJECT_ID> \
  -e API_URL=<STACK_WEBPOINT_URL> \
  -e TEST_NAME=<TEST_NAME> \
  -e BOOK_OP=<BOOK_OP> \
  ./tests/performance/main.js
```

### Run a test locally and stream the results to Grafana Cloud

```shell
./k6 run -e K6_CLOUD_PROJECT_ID=<PROJECT_ID>\
  -e API_URL=<STACK_WEBPOINT_URL> \
  -e TEST_NAME=<TEST_NAME> \
  -e BOOK_OP=<BOOK_OP=<BOOK_OP> \
  -e WORKLOAD=<WORKLOAD> \
  --out cloud \
  --out dashboard=report=<REPORT_NAME>.html \
  ./tests/performance/main.js
```
