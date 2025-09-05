# Go Modular Monolith REST API

A production-ready modular monolith REST API built with Go, featuring authentication, user management, and product catalog functionality.

## 🚀 Features

- **Modular Architecture**: Feature-based folder structure with clear separation of concerns
- **Authentication**: JWT-based auth with access/refresh tokens and bcrypt password hashing
- **User Management**: Complete CRUD operations with profile management
- **Product Catalog**: Product management with categories and search capabilities
- **Database**: PostgreSQL with GORM ORM and automatic migrations
- **Observability**: OpenTelemetry tracing, structured logging, and health checks
- **Security**: Rate limiting, CORS, input validation, and secure middleware
- **Development**: Docker support, live reload, comprehensive testing

## 🛠 Tech Stack

- **Framework**: Echo v4
- **Database**: PostgreSQL with GORM
- **Dependency Injection**: Uber FX
- **Configuration**: Viper
- **Validation**: go-playground/validator v10
- **Logging**: slog (structured logging)
- **Telemetry**: OpenTelemetry with Jaeger
- **Authentication**: JWT with golang-jwt/jwt v5
- **Password Hashing**: bcrypt
- **Migrations**: Goose
- **Task Runner**: Task

## 📁 Project Structure

```
.
├── cmd/
│   ├── server/          # HTTP server entry point
│   └── migrate/         # Database migration tool
├── internal/
│   ├── auth/            # Authentication module
│   ├── user/            # User management module
│   ├── product/         # Product management module
│   ├── database/        # Database connection and utilities
│   ├── errors/          # Error handling
│   ├── logger/          # Structured logging
│   ├── middleware/      # HTTP middleware
│   ├── telemetry/       # OpenTelemetry setup
│   ├── validator/       # Request validation
│   ├── migration/       # Database migrations
│   ├── config/          # Configuration management
│   └── app/             # Application setup and DI
├── migrations/          # Database migrations
├── config/              # Configuration files
├── Taskfile.yml         # Task automation
├── docker-compose.yml   # Production Docker setup
└── docker-compose.dev.yml # Development services
```

## 🚦 Quick Start

### Prerequisites

- Go 1.23+
- Docker and Docker Compose
- Task (task runner) - `go install github.com/go-task/task/v3/cmd/task@latest`

### Setup

1. **Clone and setup**:
   ```bash
   git clone <repository-url>
   cd go-modulith
   task dev:setup
   ```

2. **Start development environment**:
   ```bash
   task app:run
   ```

3. **API will be available at**: `http://localhost:8080`

### Alternative Docker Setup

```bash
# Start all services with Docker
task docker:run

# Or run in background
task docker:run-detached
```

## 📚 API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - Logout (requires auth)

### Users
- `GET /api/v1/users` - List users (paginated, requires auth)
- `GET /api/v1/users/profile` - Get current user profile (requires auth)
- `GET /api/v1/users/:id` - Get user by ID (requires auth)
- `PUT /api/v1/users/:id` - Update user (requires auth)
- `DELETE /api/v1/users/:id` - Delete user (soft delete, requires auth)

### Products
- `GET /api/v1/products` - List products (public, paginated, filterable)
- `GET /api/v1/products/:id` - Get product by ID (public)
- `POST /api/v1/products` - Create product (requires auth)
- `PUT /api/v1/products/:id` - Update product (requires auth, owner only)
- `DELETE /api/v1/products/:id` - Delete product (requires auth, owner only)

### Health & Monitoring
- `GET /health` - Health check
- `GET /ready` - Readiness check
- `GET /version` - Version information

## 🛠 Development

### Available Tasks

```bash
task                    # Show all available tasks
task app:run           # Run the application
task app:build         # Build binaries
task db:migrate-up     # Run migrations
task test:unit         # Run unit tests
task test:integration  # Run integration tests
task tools:lint        # Run linter
task watch:dev         # Watch mode with auto-restart
```

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_modulith

# JWT
JWT_SECRET=your-super-secret-jwt-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=24h

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# OpenTelemetry
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317
OTEL_SERVICE_NAME=go-modulith
```

### Database Operations

```bash
task db:start          # Start PostgreSQL
task db:migrate-up     # Run migrations
task db:migrate-down   # Rollback migration
task db:migrate-status # Check migration status
task db:connect        # Connect with psql
```

## 🧪 Testing

```bash
# Unit tests
task test:unit

# Integration tests (requires DB)
task test:integration

# Test coverage
task test:coverage

# Benchmarks
task test:bench
```

## 📊 Observability

### Logging
- Structured JSON logging with slog
- Request/response logging
- Distributed tracing correlation

### Monitoring
- **Jaeger UI**: http://localhost:16686
- **Health endpoint**: http://localhost:8080/health
- **Ready endpoint**: http://localhost:8080/ready

### Metrics
- OpenTelemetry traces for all requests
- Database query tracing
- Custom business logic spans

## 🔒 Security Features

- JWT authentication with refresh tokens
- bcrypt password hashing
- Rate limiting (configurable)
- CORS configuration
- Request validation
- SQL injection prevention
- Input sanitization

## 🏗 Architecture Decisions

### Modular Monolith
- Feature-based modules (auth, user, product)
- Shared infrastructure components
- Clear module boundaries
- Easy to extract to microservices later

### No DTOs/Mappers
- Direct entity-to-response mapping
- Simplified data flow
- Less boilerplate code

### Dependency Injection
- Uber FX for clean DI
- Lifecycle management
- Easy testing and mocking

## 📦 Deployment

### Docker Production

```bash
# Build and run
task docker:build
task docker:run

# Health check
task prod:health
```

### Manual Deployment

```bash
# Build binary
task app:build

# Run migrations
task db:migrate-up

# Start server
./bin/go-modulith
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`task ci:test`)
4. Commit changes (`git commit -m 'Add amazing feature'`)
5. Push to branch (`git push origin feature/amazing-feature`)
6. Create Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔧 Troubleshooting

### Common Issues

1. **Database connection failed**:
   ```bash
   task db:start  # Ensure PostgreSQL is running
   task db:logs   # Check database logs
   ```

2. **Migration errors**:
   ```bash
   task db:migrate-status  # Check current status
   task db:migrate-reset   # Reset if needed
   ```

3. **Port already in use**:
   ```bash
   lsof -i :8080  # Find process using port
   kill -9 <PID>  # Kill process
   ```

### Development Setup Issues

```bash
# Reset everything and start fresh
task dev:reset
```

## 📖 API Documentation

The API follows RESTful conventions with consistent response formats:

### Success Response
```json
{
  "id": "uuid",
  "data": {...},
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Error Response
```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {
    "field": "Validation error"
  }
}
```

### Pagination
```json
{
  "data": [...],
  "total": 100,
  "page": 1,
  "limit": 10,
  "total_pages": 10
}
```

For detailed API documentation, run the application and visit the generated docs or use tools like Postman with the provided endpoints.