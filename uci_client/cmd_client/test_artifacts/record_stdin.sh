#!/bin/bash

echo "recording stdin..."
cd "$(dirname "$0")" || exit

# Define the output file
output_file="output.txt"

# Read from stdin and write to the output file
while IFS= read -r line
do
  echo "$line" > "$output_file"
  echo "recorded: $line"
done