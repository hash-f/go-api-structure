#!/bin/bash

# Script to generate and validate Swagger documentation

# Ensure script stops on first error
set -e

echo "Generating Swagger documentation..."

# Generate Swagger docs using the swag CLI
# Generate Swagger docs using the swag CLI
# Use -g to specify the main Go file as the entry point
# Use -o to specify the output directory for generated docs
swag init -g cmd/api/main.go -o internal/docs

if [ $? -eq 0 ]; then
  echo "Swagger documentation generated successfully in internal/docs."
  echo "You can access the Swagger UI at: http://localhost:8080/swagger/index.html when the server is running"
else
  echo "Error generating Swagger documentation."
  exit 1
fi
