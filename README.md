# Doan - Go Monorepo Project(Đồ án tốt nghiệp Project BE - Nguyễn Thế Hà 61165 CS2)

Clean architecture Go monorepo with GORM, Wire DI, and Gin framework.

## 🚀 Quick Start

```bash
# Install dependencies
make install-tools
make deps

# Run migration
make migrate

# Start server
make dev
```

**Access:**
- API: http://localhost:8080/api
- Swagger: http://localhost:8080/api/swagger/index.html
- Health: http://localhost:8080/api/ping

## 📋 Prerequisites

- Go 1.25+
- PostgreSQL 14+
- Wire (auto-installed via `make install-tools`)
- Swag (auto-installed via `make install-tools`)

## ⚙️ Configuration

1. Copy config file:
```bash
cp configs/config.yaml.sample configs/config.yaml
```

2. Update database settings in `configs/config.yaml`

## 🔧 Common Commands

```bash
# Development
make dev              # Run dev server (auto-generate)
make wire             # Generate dependency injection
make migrate          # Run database migration

# Build
make build            # Build HTTP server
make build-migration  # Build migration CLI

# Testing
make test             # Run tests
make test-coverage    # Test with coverage

# Utilities
make clean            # Clean generated files
make fmt              # Format code
```

## 📁 Project Structure

```
.
├── cmd/                    # Application entrypoints
│   ├── http/              # HTTP REST API server
│   └── cli/migration/     # Database migration CLI
├── internal/              # Private application code
│   ├── entities/          # Domain entities
│   ├── repositories/      # Data access layer
│   ├── services/          # Business logic
│   ├── usecases/          # Application logic
│   └── infrastructure/    # External implementations
├── pkg/                   # Public libraries
│   ├── config/           # Configuration management
│   ├── logger/           # Logging utilities
│   └── utils/            # Helper functions
├── configs/              # Configuration files
└── scripts/              # Utility scripts
```

## 🗄️ Database Migration

Using GORM AutoMigrate:

```bash
# Run migration
make migrate
```

## 🔌 Wire Dependency Injection

Dependencies are auto-wired via Google Wire:

```bash
# Generate wire code
make wire

# Files generated:
# - cmd/http/wire_gen.go
# - cmd/cli/migration/wire_gen.go
```

## 📖 API Documentation

Swagger documentation auto-generated: //todo

```bash
# Generate swagger docs
make swagger

# Access: http://localhost:8080/api/swagger/index.html
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage
open coverage.html
```

## 🐳 Docker

```bash
# Build image
make docker-build

# Start local services (postgres, redis, etc.)
make docker-local-up

# Stop services
make docker-local-down
```

## 📝 Adding New Entity

1. Create entity in `internal/entities/`:
```go
type Student struct {
    ID   uuid.UUID `gorm:"type:uuid;primary_key"`
    Name string    `gorm:"not null"`
}
```

2. Add to migration in `internal/infrastructure/database/postgres/migration.go`:
```go
func (m *migration) getAllEntities() []interface{} {
    return []interface{}{
        &entities.User{},
        &entities.Student{},  // Add here
    }
}
```

3. Run migration:
```bash
make migrate
```

## 🛠️ Troubleshooting

### Wire generation fails
```bash
make clean
make install-tools
make wire
```

### Port already in use
```bash
# Kill process on port 8080
lsof -ti:8080 | xargs kill -9
```

## 📚 Documentation

- **START_HERE.md** - Quick start guide
- **SETUP.md** - Detailed setup instructions  
- **DEVELOPMENT.md** - Development workflow
- **QUICK_REFERENCE.md** - Command reference

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👥 Authors

- Nguyễn Thế Hà 61165 CS2 - Initial work(đồ án tốt nghiệp Back-end project)

## 🙏 Acknowledgments

- [GORM](https://gorm.io/) - ORM library
- [Wire](https://github.com/google/wire) - Dependency injection
- [Gin](https://github.com/gin-gonic/gin) - Web framework
- [Zap](https://github.com/uber-go/zap) - Logging library

