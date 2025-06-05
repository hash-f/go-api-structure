#!/bin/bash

# Script to generate and validate Swagger documentation

# Ensure script stops on first error
set -e

echo "Generating Swagger documentation..."

# Generate Swagger docs using the swag CLI
~/go/bin/swag init \
  --dir ./cmd/api,./internal/api \
  --parseDependency \
  --output ./internal/docs \
  --generalInfo ./internal/docs/swagger.go \
  --parseInternal

echo "Swagger documentation generated successfully!"
echo "You can access the Swagger UI at: http://localhost:8080/swagger/index.html when the server is running"
