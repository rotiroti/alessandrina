#!/usr/bin/env bash
#
# Usage: ./k6deployment.sh <STACK_NAME> <AWS_REGION> <WORKLOAD> <K6_SCRIPT>

# Output directory for the reports
REPORTS_DIR=$(pwd)/reports

# K6 binary path
K6_BINARY=$(pwd)/k6

# Supported k6 workloads
WORKLOADS=("baseline" "vus5" "vus10" "vus15" "averageVUs" "stressVUs" "stressRate")

# Check if the stack name is provided
if [ $# -eq 0 ]; then
    echo "Please provide the stack name"
    exit 1
fi

stack_name=$1

# Check if the AWS region is provided
if [ $# -eq 1 ]; then
    echo "Please provide the AWS region"
    exit 1
fi

aws_region=$2

# Check if the k6 workload is provided
if [ $# -eq 2 ]; then
    echo "Please provide the k6 workload"
    exit 1
fi

workload_idx=$3

# Check if the K6 script is provided
if [ $# -eq 3 ]; then
    echo "Please provide the K6 script"
    exit 1
fi

k6_script=$4

# Check if the provided workload index is within the valid range
if [ "$workload_idx" -ge 4 ] && [ "$workload_idx" -lt ${#WORKLOADS[@]} ]; then
    workload=${WORKLOADS[$workload_idx]}
else
    echo "Invalid workload index. Possible values: 0-${#WORKLOADS[@]}"
    exit 1
fi

# Retrieve the API Gateway URL from the stack outputs.
api_url=$(sam list stack-outputs --stack-name "$stack_name" --region "$aws_region" --output json | jq -r '.[] | select(.OutputKey=="WebEndpoint")|.OutputValue')

# Check if the API Gateway URL is set.
if [ -z "$api_url" ]; then
    echo "API Gateway URL not found"
    exit 1
fi

# Check if the REPORTS_DIR directory exists.
if [ ! -d "$REPORTS_DIR" ]; then
    mkdir -p "$REPORTS_DIR"
fi

# Define the output directory for the k6 workload summary reports.
OUTPUT_DIR="$REPORTS_DIR/$stack_name/$workload/$(date +"%Y-%m-%d-%H-%M-%S")/"

# Check if the OUTPUT_DIR directory exists.
if [ ! -d "$OUTPUT_DIR" ]; then
    mkdir -p "$OUTPUT_DIR"
fi

$K6_BINARY run -e WORKLOAD="$workload_idx" -e API_URL="$api_url" --out dashboard=report="$OUTPUT_DIR/report.html" --summary-trend-stats="avg,min,med,max,p(90),p(95),p(99)" "$k6_script"