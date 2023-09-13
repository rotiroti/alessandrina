#!/usr/bin/env bash

# Generate synthetic data using k6
# Usage: ./stack-seed.sh <STACK_NAME> <AWS_REGION>

N=10
OUTPUT_DIR=$(pwd)/reports/seeds

K6_BINARY=$(pwd)/k6
K6_SCRIPT=$(pwd)/tests/performance/create-book.js

# Check if the stack name is provided
if [ $# -eq 0 ]; then
    echo "Please provide the stack name"
    exit 1
fi

# Check if the AWS region is provided
if [ $# -eq 1 ]; then
    echo "Please provide the AWS region"
    exit 1
fi

stack_name=$1
aws_region=$2

# Retrieve the API Gateway URL from the stack outputs.
api_url=$(sam list stack-outputs --stack-name "$stack_name" --region "$aws_region" --output json | jq -r '.[] | select(.OutputKey=="WebEndpoint")|.OutputValue')

# Check if the API Gateway URL is set.
if [ -z "$api_url" ]; then
    echo "API Gateway URL not found"
    exit 1
fi

for ((i=1; i <= N; i++));
do
    # Run the baseline workload using the create-book.js script and save the results in a HTML file.
    $K6_BINARY run -e API_URL="$api_url" --out dashboard=report="$OUTPUT_DIR/$stack_name/$(date +"%Y-%m-%d-%H-%M-%S")/report.$i.html" --summary-trend-stats="avg,min,med,max,p(90),p(95),p(99)" "$K6_SCRIPT";
done
