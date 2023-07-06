#!/usr/bin/env bash

# Delete a DynamoDB table using the AWS CLI and the localstack endpoint
# Usage: ./delete-table.sh <table_name>

# Check if the table name is provided
if [ $# -eq 0 ]; then
    echo "Please provide the table name"
    exit 1
fi

# Get the table name
table_name=$1

# Check if the table exists
if ! aws dynamodb describe-table \
    --endpoint-url http://localhost:4566 \
    --table-name "$table_name" > /dev/null 2>&1; then
    echo "Table $table_name does not exists"
    exit 1
fi

# Delete the table
aws dynamodb delete-table \
    --endpoint-url http://localhost:4566 \
    --table-name "$table_name"
