#!/bin/bash

# Development script for running the Go server with hot reload

echo "Starting development server with hot reload..."
echo "Server will restart automatically when you save changes to .go files"
echo ""

# Set up Go bin path
export PATH="$HOME/go/bin:$PATH"

# Check if air is installed
if ! command -v air &> /dev/null; then
    echo "Air is not installed. Installing..."
    go install github.com/air-verse/air@latest
fi

# Run air
air