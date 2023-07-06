#!/usr/bin/env bash

# Create a DynamoDB table using the AWS CLI and the localstack endpoint
# Usage: ./create-table.sh <table_name>

# Check if the table name is provided
if [ $# -eq 0 ]; then
    echo "Please provide the table name"
    exit 1
fi

# Get the table name
table_name=$1

# Check if the table exists
if aws dynamodb describe-table \
    --endpoint-url http://localhost:4566 \
    --table-name "$table_name" > /dev/null 2>&1; then
    echo "Table $table_name already exists"
    exit 1
fi

# Create the table
aws dynamodb create-table \
    --endpoint-url http://localhost:4566 \
    --table-name "$table_name" \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1
