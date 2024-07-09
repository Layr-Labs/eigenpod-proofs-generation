#!/bin/bash

# Check if the input file is provided
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <input_file>"
    exit 1
fi

input_file=$1
temp_file=$(mktemp)

# Check if jq is installed
if ! command -v jq &> /dev/null; then
    echo "jq is required but not installed. Please install jq (brew install jq)."
    exit 1
fi

# Extract the abi property and overwrite the file
jq '.abi' "$input_file" > temp_file && mv temp_file "$input_file"

echo "File has been updated with the abi content at the root level."