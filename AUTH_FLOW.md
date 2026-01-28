se# Authentication Flow - Complete Guide

## Overview
This document describes the complete authentication flow implemented in this project using JWT tokens, following clean architecture principles with GORM and dependency injection via Wire.

## Architecture

```
┌─────────────┐
│  Controller │  <- HTTP Layer (Gin)
└──────┬──────┘
       │
┌──────▼──────┐
│   UseCase   │  <- Business Logic
└──────┬──────┘
       │
┌──────▼──────┐
│   Service   │  <- Domain Services
└──────┬──────┘
       │
┌──────▼──────┐
│ Repository  │  <- Data Access (GORM)
└──────┬──────┘
       │
┌──────▼──────┐
│  Database   │  <- PostgreSQL
└─────────────┘
```

## Authentication Endpoints

### 1. Login
**POST** `/v1/auth/login`

**Request:**
```json
{
  "username": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "code": "USER001",
      "full_name": "John Doe",
      "email": "user@example.com",
      "role": "STUDENT",
      "is_active": true
    }
  }
}
```

### 2. Logout
**POST** `/v1/auth/logout`

**Request:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Logout successful",
  "data": {
    "message": "Logged out successfully"
  }
}
```

### 3. Refresh Token
**POST** `/v1/auth/refresh`

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## Project Structure

```
.
├── cmd/
│   ├── http/                      # HTTP server
│   │   ├── main.go               # Entry point
│   │   ├── wire.go               # Wire dependency injection
│   │   ├── wire_gen.go           # Generated wire code
│   │   ├── controllers/          # HTTP handlers
│   │   │   └── user/
│   │   │       ├── controller.go # Interface & routes
│   │   │       ├── v1.go         # V1 implementation
│   │   │       ├── v2.go         # V2 implementation
│   │   │       └── dto.go        # Request/Response DTOs
│   │   └── middleware/           # HTTP middleware
│   │       ├── auth.go           # JWT authentication
│   │       ├── cors.go           # CORS configuration
│   │       └── ginLogger.go      # Logging middleware
│   └── cli/
│       └── migration/            # Database migration CLI
│           ├── main.go
│           └── wire.go
├── internal/
│   ├── entities/                 # Domain entities (GORM models)
│   │   └── user.go
│   ├── repositories/             # Data access layer
│   │   ├── interface/
│   │   │   └── user.go          # Repository interface
│   │   └── base.go              # Base repository
│   ├── services/                 # Domain services
│   │   └── user/
│   │       ├── auth.go          # Authentication service
│   │       └── dto.go           # Service DTOs
│   ├── usecases/                # Business logic
│   │   └── user/
│   │       ├── login.go         # Login usecase
│   │       ├── logout.go        # Logout usecase
│   │       └── refresh_token.go # Refresh token usecase
│   └── infrastructure/          # External dependencies
│       └── database/
│           └── postgres/
│               ├── gorm.go      # GORM setup
│               └── migration.go # Auto migration
├── pkg/                         # Shared packages
│   ├── utils/
│   │   ├── jwt.go              # JWT utilities
│   │   └── auth.go             # Auth utilities
│   ├── config/                 # Configuration management
│   ├── logger/                 # Logging utilities
│   └── types/                  # Common types
└── configs/
    └── config.yaml             # Application configuration
```

## Setup Instructions

### 1. Prerequisites
- Go 1.25+
- PostgreSQL 14+
- Make

### 2. Install Dependencies
```bash
# Install required tools
make install-tools

# Install Go dependencies
go mod tidy
```

### 3. Configure Application

Edit `configs/config.yaml`:
```yaml
http:
  path: 0.0.0.0
  port: 9000

database:
  driver: postgres
  user: root
  password: password
  protocol: tcp
  host: localhost
  port: 5432
  schema: doan

redis:
  url: localhost:6379

jwt:
  secret: "your-secret-key-change-this-in-production"
  access_token_duration: 24h
  refresh_token_duration: 168h
```

### 4. Setup Database

```bash
# Create PostgreSQL database
createdb doan

# Or using psql
psql -U postgres -c "CREATE DATABASE doan;"
```

### 5. Run Migrations

```bash
# Generate wire dependencies and run migration
make migrate
```

This will:
- Generate wire dependency injection code
- Create database tables using GORM AutoMigrate
- Enable UUID extension in PostgreSQL

### 6. Generate Code

```bash
# Generate Wire DI code and Swagger docs
make generate
```

This will:
- Generate `wire_gen.go` files
- Generate Swagger documentation

### 7. Run Application

```bash
# Development mode (with hot reload support)
make dev

# Or standard run
make run
```

The server will start on `http://localhost:9000`

### 8. Access Swagger Documentation

Open your browser and navigate to:
```
http://localhost:9000/api/swagger/index.html
```

## Authentication Flow Details

### Login Flow

1. **Client** sends POST request to `/v1/auth/login` with username and password
2. **Controller** validates request body and calls LoginUseCase
3. **LoginUseCase** calls AuthService.CreateAuthToken()
4. **AuthService**:
   - Queries user from database via UserRepository
   - Validates user status (is_active)
   - Verifies password using bcrypt
   - Generates access token (24h expiry)
   - Generates refresh token (7 days expiry)
   - Returns tokens and user info
5. **Controller** returns response with tokens and user data

### Logout Flow

1. **Client** sends POST request to `/v1/auth/logout` with token
2. **Controller** validates request and calls LogoutUseCase
3. **LogoutUseCase** calls AuthService.ValidateToken()
4. **AuthService** validates the token
5. Token is invalidated (in production, add to Redis blacklist)
6. **Controller** returns success message

### Refresh Token Flow

1. **Client** sends POST request to `/v1/auth/refresh` with refresh_token
2. **Controller** validates request and calls RefreshTokenUseCase
3. **RefreshTokenUseCase** calls AuthService.RefreshAccessToken()
4. **AuthService**:
   - Validates refresh token
   - Extracts user claims
   - Generates new access token
5. **Controller** returns new access token

### Protected Routes

To protect routes with JWT authentication:

```go
import "doan/cmd/http/middleware"

// In your route registration
protected := router.Group("/api")
protected.Use(middleware.AuthMiddleware(configManager))
{
    protected.GET("/profile", controller.GetProfile)
    protected.PUT("/profile", controller.UpdateProfile)
}

// With role-based access control
admin := router.Group("/admin")
admin.Use(
    middleware.AuthMiddleware(configManager),
    middleware.RoleMiddleware("ADMIN", "TEACHER"),
)
{
    admin.GET("/users", controller.ListUsers)
}
```

## Dependency Injection with Wire

### Wire Providers

**Database Provider** (`internal/infrastructure/database/provider.go`):
```go
var DBProvider = wire.NewSet(
    GetDBContext,
    implement.NewUserRepository,
)
```

**Service Provider** (`internal/services/provider.go`):
```go
var UserServiceProvider = wire.NewSet(
    user.NewAuthService,
)
```

**UseCase Provider** (`internal/usecases/provider.go`):
```go
var UserUseCaseProviders = wire.NewSet(
    user.NewLoginUseCase,
    user.NewLogoutUseCase,
    user.NewRefreshTokenUseCase,
)
```

**Controller Provider** (`cmd/http/controllers/provider.go`):
```go
var ControllerProviders = wire.NewSet(
    user.NewUserControllerV1,
    user.NewUserControllerV2,
)
```

### Regenerate Wire Code

Whenever you modify constructors or add new dependencies:

```bash
make wire
```

## Makefile Commands

```bash
# Install development tools
make install-tools

# Generate Wire DI code
make wire

# Generate Swagger documentation
make swagger

# Generate all (wire + swagger)
make generate

# Run migrations
make migrate

# Run HTTP server
make run

# Run in development mode
make dev

# Build binary
make build

# Run tests
make test

# Clean generated files
make clean
```

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE,
    full_name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    password TEXT,
    role VARCHAR(50) DEFAULT 'STUDENT',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_code ON users(code);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

## Security Best Practices

1. **JWT Secret**: Change the JWT secret in production to a strong, random string
2. **Password Hashing**: Passwords are hashed using bcrypt
3. **Token Expiry**: Access tokens expire in 24h, refresh tokens in 7 days
4. **HTTPS**: Always use HTTPS in production
5. **Token Blacklist**: Implement Redis-based token blacklist for logout
6. **Rate Limiting**: Add rate limiting middleware for auth endpoints
7. **Input Validation**: All inputs are validated using Gin's binding

## Testing Authentication

### Using cURL

**Login:**
```bash
curl -X POST http://localhost:9000/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "user@example.com",
    "password": "password123"
  }'
```

**Access Protected Route:**
```bash
curl -X GET http://localhost:9000/api/profile \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

**Refresh Token:**
```bash
curl -X POST http://localhost:9000/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

**Logout:**
```bash
curl -X POST http://localhost:9000/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "token": "YOUR_ACCESS_TOKEN"
  }'
```

## Troubleshooting

### Wire Generation Fails

```bash
# Clean and regenerate
rm -f cmd/http/wire_gen.go cmd/cli/migration/wire_gen.go
make wire
```

### Migration Fails

```bash
# Check database connection
psql -U root -d doan -h localhost -p 5432

# Enable UUID extension manually
psql -U root -d doan -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"

# Retry migration
make migrate
```

### JWT Token Invalid

- Check JWT secret matches in config.yaml
- Verify token hasn't expired
- Ensure token format is: `Bearer <token>`

## Next Steps

1. **Add More Entities**: Create repositories, services, and usecases for other entities
2. **Implement Token Blacklist**: Use Redis to store invalidated tokens
3. **Add Email Verification**: Implement email verification flow
4. **Password Reset**: Add forgot password functionality
5. **2FA**: Implement two-factor authentication
6. **OAuth**: Add social login (Google, Facebook, etc.)
7. **Refresh Token Rotation**: Implement refresh token rotation for better security

## Contributing

When adding new features:

1. Create entity in `internal/entities/`
2. Create repository interface in `internal/repositories/interface/`
3. Implement repository in `internal/infrastructure/database/postgres/implement/`
4. Create service in `internal/services/`
5. Create usecases in `internal/usecases/`
6. Create controllers in `cmd/http/controllers/`
7. Add providers to respective `provider.go` files
8. Run `make generate` to regenerate wire and swagger docs
9. Update this README with new endpoints

## License

[Your License Here]

