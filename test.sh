#!/bin/bash

echo "Running API tests for HTML to Go conversion..."
go test -v

# If we want to measure test coverage as well
# echo "Running tests with coverage..."
# go test -v -coverprofile=coverage.out
# go tool cover -html=coverage.out -o coverage.html
# echo "Coverage report generated at coverage.html"
