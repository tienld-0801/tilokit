#!/bin/bash

# TiLoKit Test Script
set -e

echo "🧪 Running TiLoKit tests..."

# Run Go tests
echo "📋 Running unit tests..."
go test -v ./...

# Run linting
echo "🔍 Running linter..."
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
else
    echo "⚠️  golangci-lint not found, skipping linting"
fi

# Run formatting check
echo "📝 Checking code formatting..."
if [ "$(gofmt -l . | wc -l)" -gt 0 ]; then
    echo "❌ Code is not formatted properly"
    gofmt -l .
    exit 1
else
    echo "✅ Code is properly formatted"
fi

# Run security check
echo "🔒 Running security check..."
if command -v gosec &> /dev/null; then
    gosec ./...
else
    echo "⚠️  gosec not found, skipping security check"
fi

echo "✅ All tests passed!"
