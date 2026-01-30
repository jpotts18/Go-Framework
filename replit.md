# GoLaravel - A Laravel-Inspired Go Framework

## Overview
GoLaravel is a Go framework that mimics Laravel's elegant design patterns and developer experience. It provides familiar Laravel concepts implemented in idiomatic Go.

## Project Structure
```
golaravel/
├── app/
│   ├── application.go         # Main application bootstrap
│   ├── container/             # Service container (dependency injection)
│   │   └── container.go
│   ├── config/                # Configuration management
│   │   └── config.go
│   ├── http/
│   │   ├── controllers/       # Base controller with validation
│   │   │   └── controller.go
│   │   ├── middleware/        # Built-in middleware (CORS, Auth, Rate Limiting, etc.)
│   │   │   └── middleware.go
│   │   ├── request/           # HTTP request wrapper
│   │   │   └── request.go
│   │   └── response/          # HTTP response wrapper
│   │       └── response.go
│   ├── routing/               # Router with groups, params, middleware
│   │   └── router.go
│   ├── database/
│   │   ├── orm/               # Query builder and ORM
│   │   │   └── model.go
│   │   └── migration/         # Schema builder and migrations
│   │       └── migration.go
│   ├── validation/            # Laravel-style validation
│   │   └── validation.go
│   ├── view/                  # Template engine
│   │   └── view.go
│   └── support/               # Helper functions (Str, Collection, Carbon, etc.)
│       └── helpers.go
└── examples/
    └── demo/                  # Example application
        └── main.go
```

## Key Features

### 1. Service Container
- Dependency injection with bindings and singletons
- Automatic resolution with `inject` struct tags
- Laravel-like `Make()` and `Instance()` methods

### 2. Routing
- Expressive route definitions: `router.Get()`, `router.Post()`, etc.
- Route parameters: `/users/{id}`
- Route groups with prefixes and middleware
- Named routes
- 404 handler customization

### 3. Middleware
- Logger, Recovery, CORS, Auth, Rate Limiting
- Secure headers, Basic Auth, Timeout
- Easy to create custom middleware

### 4. Request/Response
- Request input helpers: `Param()`, `Query()`, `JSON()`, `All()`
- Response methods: `JSON()`, `HTML()`, `Redirect()`
- Error responses: `NotFound()`, `BadRequest()`, `ValidationError()`

### 5. Validation
- Laravel-style rules: `required`, `email`, `min`, `max`, `between`, `in`, etc.
- Easy integration with controllers

### 6. Database
- Query builder with fluent interface
- ORM with models
- Migration system with schema builder

### 7. Views
- Template engine with shared data
- Built-in helper functions

## Running the Demo
```bash
cd golaravel
go run ./examples/demo/
```

The server starts on port 5000.

## Available Demo Endpoints
- `GET /` - Home page
- `GET /about` - About information
- `GET /user/{id}` - Get user by ID
- `POST /users` - Create user (with validation)
- `GET /api/status` - API health status
- `GET /api/v1/products` - List products with pagination
- `GET /validate?email=...&name=...&age=...&password=...` - Validation demo

## Running Tests
```bash
cd golaravel
go test ./...
```

All 326 tests pass, covering:
- Service container (bindings, singletons, instances)
- Routing (groups, params, middleware, named routes)
- Request/Response handling
- Validation (all Laravel-style rules)
- Config (nested keys, environment loading)
- Views (templates, shared data)
- Middleware (CORS, auth, rate limiting, etc.)
- Database ORM (query builder, models)
- Migrations (schema builder)
- Support helpers (Str, Collection, Carbon)
- Authentication (password hashing, token guards)
- Caching (in-memory store with expiration)
- Sessions (storage, flash data, regeneration)
- Events (dispatcher, listeners, subscribers)

## Documentation Website
The library includes a hosted documentation site with:
- Homepage with features overview
- Comprehensive documentation covering all framework components
- Sidebar navigation for easy browsing

Run the docs server:
```bash
cd golaravel && go run ./docs/
```

## Recent Changes
- Created initial framework structure (2026-01-30)
- Implemented all core components
- Added example demo application
- Added comprehensive test suite with 273 tests (2026-01-30)
- Fixed nullable validation to skip subsequent rules when field is empty
- Fixed Config.Get() to support flat keys from LoadEnv
- Added documentation website with comprehensive guides (2026-01-30)
- Added transaction support for database operations (2026-01-30)
- Added model relationships: HasOne, HasMany, BelongsTo, BelongsToMany (2026-01-30)
- Added auth package with bcrypt password hashing (2026-01-30)
- Added cache package with in-memory store (2026-01-30)
- Added session package with flash data support (2026-01-30)
- Added events package with dispatcher and subscribers (2026-01-30)
- Added graceful shutdown handling (2026-01-30)
- Added context.Context support to Request (2026-01-30)
- Fixed rate limiter concurrency with sync.Mutex (2026-01-30)
