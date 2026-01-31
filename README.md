# GoLaravel

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A Laravel-inspired web framework for Go that brings the elegance and developer experience of Laravel to the Go ecosystem.

## Why GoLaravel?

If you love Laravel's expressive syntax and comprehensive tooling but need Go's performance and type safety, GoLaravel bridges that gap. It implements familiar Laravel patterns in idiomatic Go, making it easy for Laravel developers to be productive in Go immediately.

## Features

- **Service Container** - Dependency injection with automatic resolution
- **Expressive Routing** - Route groups, parameters, and middleware support
- **Middleware Stack** - CORS, Auth, Rate Limiting, and more built-in
- **Validation** - Laravel-style validation rules (`required`, `email`, `min`, `max`, etc.)
- **ORM & Query Builder** - Fluent database interface with relationships
- **Migrations** - Schema builder for database versioning
- **Template Engine** - Views with shared data and helper functions
- **Authentication** - Password hashing and token guards
- **Caching** - In-memory store with expiration
- **Sessions** - Flash data and session management
- **Events** - Event dispatcher with listeners and subscribers

## Quick Start

### Installation

```bash
go get github.com/yourusername/golaravel
```

### Hello World

```go
package main

import (
    "golaravel/app"
    "golaravel/app/http/response"
)

func main() {
    application := app.New()

    application.Router().Get("/", func(res *response.Response) {
        res.JSON(200, map[string]string{
            "message": "Hello, GoLaravel!",
        })
    })

    application.Run(":5000")
}
```

### Running the Demo

```bash
cd golaravel
go run ./examples/demo/
```

The server starts on port 5000 with these endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Home page |
| GET | `/about` | About information |
| GET | `/user/{id}` | Get user by ID |
| POST | `/users` | Create user (with validation) |
| GET | `/api/status` | API health status |
| GET | `/api/v1/products` | List products with pagination |
| GET | `/validate` | Validation demo |

## Core Concepts

### Routing

```go
router := app.Router()

// Basic routes
router.Get("/users", listUsers)
router.Post("/users", createUser)
router.Put("/users/{id}", updateUser)
router.Delete("/users/{id}", deleteUser)

// Route groups with middleware
router.Group("/api", func(r *routing.Router) {
    r.Use(middleware.Auth())
    r.Get("/profile", getProfile)
})
```

### Middleware

```go
// Built-in middleware
router.Use(middleware.Logger())
router.Use(middleware.Recovery())
router.Use(middleware.CORS())
router.Use(middleware.RateLimiter(100, time.Minute))

// Custom middleware
func MyMiddleware() routing.Middleware {
    return func(next routing.Handler) routing.Handler {
        return func(req *request.Request, res *response.Response) {
            // Before request
            next(req, res)
            // After request
        }
    }
}
```

### Validation

```go
rules := map[string]string{
    "email":    "required|email",
    "name":     "required|min:2|max:50",
    "age":      "required|numeric|between:18,120",
    "password": "required|min:8",
}

errors := validator.Validate(data, rules)
if len(errors) > 0 {
    res.ValidationError(errors)
    return
}
```

### Database ORM

```go
// Query builder
users := db.Table("users").
    Where("active", true).
    OrderBy("created_at", "desc").
    Limit(10).
    Get()

// Model relationships
type User struct {
    orm.Model
    Posts []Post `rel:"hasMany"`
}

user.HasMany(&user.Posts)
```

### Service Container

```go
container := container.New()

// Register a singleton
container.Singleton("db", func() interface{} {
    return database.Connect()
})

// Resolve dependencies
db := container.Make("db").(*database.Connection)
```

## Project Structure

```
golaravel/
├── app/
│   ├── application.go          # Main application bootstrap
│   ├── container/              # Dependency injection
│   ├── config/                 # Configuration management
│   ├── http/
│   │   ├── controllers/        # Base controller
│   │   ├── middleware/         # Built-in middleware
│   │   ├── request/            # HTTP request wrapper
│   │   └── response/           # HTTP response wrapper
│   ├── routing/                # Router implementation
│   ├── database/
│   │   ├── orm/                # Query builder and ORM
│   │   └── migration/          # Schema builder
│   ├── validation/             # Validation engine
│   ├── view/                   # Template engine
│   └── support/                # Helper functions
├── examples/
│   └── demo/                   # Example application
└── docs/                       # Documentation site
```

## Testing

```bash
go test ./...
```

The test suite includes 327 tests covering all framework components.

## Documentation

Run the documentation server locally:

```bash
cd golaravel && go run ./docs/
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [Laravel](https://laravel.com/) - The PHP Framework for Web Artisans
- Built with love for the Go community
