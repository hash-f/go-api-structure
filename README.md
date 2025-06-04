# Go API Structure

A production-ready Go API starter repository with a clean architecture design, implementing best practices for building scalable, maintainable web services.

## Overview

This repository provides a structured foundation for building robust REST APIs in Go. It follows clean architecture principles to ensure separation of concerns, testability, and maintainability.

Key features:

- Structured project layout with domain-driven design
- Configuration management with environment variables
- PostgreSQL integration with migrations
- Chi router with middleware setup
- JWT authentication
- Structured logging with slog
- Graceful shutdown handling

## Architecture

The project follows a layered architecture:

```
go-api-structure/
├── cmd/api/           # Application entrypoint
├── internal/          # Private application packages
│   ├── api/           # HTTP handlers and API utilities
│   ├── auth/          # Authentication logic
│   ├── config/        # Configuration management
│   ├── database/      # Database connection management
│   ├── logger/        # Logging setup
│   ├── server/        # HTTP server implementation
│   └── store/         # Data access layer
├── migrations/        # Database migration files
└── sqlc.yaml          # SQLC configuration
```

### Design Principles

1. **Separation of Concerns**: Each component has a specific responsibility
2. **Dependency Injection**: Dependencies are passed explicitly, making testing easier
3. **Interface-Driven Design**: Core business logic depends on interfaces, not implementations
4. **Error Handling**: Consistent error handling throughout the application
5. **Context Propagation**: Context is used for cancellation, timeouts, and tracing

## Tech Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL
- **Web Framework**: Chi router
- **SQL Generation**: SQLC for type-safe database access
- **Authentication**: JWT with bcrypt password hashing
- **Logging**: slog structured logging
- **Migration**: SQL migration files

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- SQLC CLI
- Make (optional, for using Makefile shortcuts)

### Database Setup

1. Create a PostgreSQL database:

```bash
createdb go_api_db
```

2. Run migrations:

```bash
# Install migrate CLI if needed
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgres://localhost:5432/go_api_db?sslmode=disable" up
```

### Configuration

Create a `.env` file in the project root with the following variables:

```
APP_ENV=development
HTTP_PORT=8080
DATABASE_DSN=postgres://postgres:postgres@localhost:5432/go_api_db?sslmode=disable
JWT_SECRET=your_very_secure_jwt_secret_key
JWT_EXPIRY_DURATION=24h
```

### Running the Application

```bash
go run cmd/api/main.go
```

## Adding New Resources

### 1. Adding a New Entity

To add a new entity (e.g., Product):

1. **Create SQL Migration**:

   ```sql
   -- in migrations/YYYYMMDDHHMMSS_create_products_table.up.sql
   CREATE TABLE products (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     name TEXT NOT NULL,
     description TEXT,
     price INTEGER NOT NULL,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
     updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
   );

   -- in migrations/YYYYMMDDHHMMSS_create_products_table.down.sql
   DROP TABLE IF EXISTS products;
   ```

2. **Define SQLC Queries**:

   ```sql
   -- in internal/store/query/product.sql
   -- name: CreateProduct :one
   INSERT INTO products (name, description, price)
   VALUES ($1, $2, $3)
   RETURNING *;

   -- name: GetProduct :one
   SELECT * FROM products WHERE id = $1;

   -- name: ListProducts :many
   SELECT * FROM products
   ORDER BY created_at DESC
   LIMIT $1 OFFSET $2;
   ```

3. **Generate SQLC Code**:

   ```bash
   sqlc generate
   ```

4. **Create Store Interface**:

   ```go
   // in internal/store/product.go
   package store

   import (
     "context"
     "go-api-structure/internal/store/db"
   )

   type ProductStore interface {
     CreateProduct(ctx context.Context, params db.CreateProductParams) (db.Product, error)
     GetProduct(ctx context.Context, id uuid.UUID) (db.Product, error)
     ListProducts(ctx context.Context, params db.ListProductsParams) ([]db.Product, error)
   }
   ```

5. **Add to Store Interface**:
   ```go
   // in internal/store/store.go
   type Store interface {
     UserStore
     ProductStore
     // other stores...
   }
   ```

### 2. Adding a New Route

1. **Create DTOs**:

   ```go
   // in internal/api/dto/product.go
   package dto

   type CreateProductRequest struct {
     Name        string  `json:"name" validate:"required"`
     Description string  `json:"description"`
     Price       int     `json:"price" validate:"required,gt=0"`
   }

   type ProductResponse struct {
     ID          string  `json:"id"`
     Name        string  `json:"name"`
     Description string  `json:"description"`
     Price       int     `json:"price"`
     CreatedAt   string  `json:"created_at"`
     UpdatedAt   string  `json:"updated_at"`
   }
   ```

2. **Create Handler**:

   ```go
   // in internal/api/handler_product.go
   package api

   // ProductHandler handles HTTP requests related to products
   type ProductHandler struct {
     store store.ProductStore
   }

   // NewProductHandler creates a new product handler
   func NewProductHandler(store store.ProductStore) *ProductHandler {
     return &ProductHandler{
       store: store,
     }
   }

   // CreateProduct handles creating a new product
   func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
     // Implementation...
   }

   // GetProduct handles retrieving a product by ID
   func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
     // Implementation...
   }

   // ListProducts handles listing products with pagination
   func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
     // Implementation...
   }
   ```

3. **Update Server**:

   ```go
   // in internal/server/server.go
   type Server struct {
     // Existing fields...
     productHandler *api.ProductHandler
   }

   func NewServer(cfg *config.Config, logger *slog.Logger, store store.Store) http.Handler {
     s := &Server{
       // Existing initialization...
       productHandler: api.NewProductHandler(store),
     }
     // Rest of the function...
   }

   // in addRoutes method, add:
   s.router.Route("/api/v1/products", func(r chi.Router) {
     r.Use(s.authService.Middleware(api.ErrorResponse))
     r.Post("/", s.productHandler.CreateProduct)
     r.Get("/{productID}", s.productHandler.GetProduct)
     r.Get("/", s.productHandler.ListProducts)
   })
   ```

### 3. Best Practices

- Keep entity logic contained in dedicated files/packages
- Follow the existing patterns for consistency
- Add validation for all incoming data
- Maintain proper error handling and logging
- Write tests for new components
