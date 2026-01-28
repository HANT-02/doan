# Monorepo Go Project - Quick Reference

## 🚀 Một Câu Lệnh Duy Nhất

```bash
# Development (Generate tất cả và chạy)
make dev
```

## 📋 Các Lệnh Thường Dùng

### Khởi Tạo (Lần Đầu)
```bash
make install-tools  # Cài Wire & Swag
make deps          # Download dependencies
make docker-local-up  # Start PostgreSQL, Redis
cp configs/config.yaml.sample configs/config.yaml
make dev           # Chạy app
```

### Development Hàng Ngày
```bash
make dev           # Generate & run
make run           # Generate & run
make run-no-gen    # Run without regenerate
```

### Code Generation
```bash
make generate      # Generate tất cả (Wire + Swagger)
make wire         # Generate Wire only
make swagger      # Generate Swagger only
```

### Build
```bash
make build         # Build HTTP server
make build-all     # Build all binaries
```

### Testing
```bash
make test          # Run tests
make test-coverage # Test với coverage report
```

### Database
```bash
make migrate       # Run migrations
```

### Docker
```bash
make docker-local-up    # Start local services
make docker-local-down  # Stop local services
make docker-build      # Build Docker image
```

### Cleanup
```bash
make clean         # Clean generated files
```

## 🎯 Workflow

### Thêm Feature Mới

1. **Entity** → `internal/entities/new_entity.go`
2. **Repository Interface** → `internal/repositories/new_repository.go`
3. **Repository Impl** → `internal/infrastructure/database/postgres/new_repository.go`
4. **Service** → `internal/services/new_service/service.go`
5. **Use Case** → `internal/usecases/new_usecase/usecase.go`
6. **Controller** → `cmd/http/controllers/new_controller/controller.go`
7. **Provider** → Update các file `provider.go`
8. **Wire** → Update `cmd/http/wire.go`
9. **Generate** → `make generate`
10. **Run** → `make dev`

### Thêm Dependency Injection

```go
// 1. Tạo provider function
func NewMyService(dep Dependency) *MyService {
    return &MyService{dep: dep}
}

// 2. Thêm vào provider set
var MyServiceProvider = wire.NewSet(NewMyService)

// 3. Update wire.go
wire.Build(
    // ...existing providers...
    MyServiceProvider,
    // ...
)

// 4. Generate
make generate
```

## 📁 Cấu Trúc Quan Trọng

```
cmd/http/
  ├── main.go           # Entry point
  ├── wire.go           # Wire config (EDIT THIS)
  └── wire_gen.go       # Generated (DON'T EDIT)

internal/
  ├── entities/         # Domain models
  ├── repositories/     # Interfaces
  ├── services/         # Business logic
  ├── usecases/         # Application logic
  └── infrastructure/   # Implementations

**/provider.go          # Wire providers
```

## 🔧 Aliases (Optional)

Thêm vào `~/.zshrc`:

```bash
alias dev="cd /Users/hant/golang/doan && make dev"
alias dgen="cd /Users/hant/golang/doan && make generate"
alias dbuild="cd /Users/hant/golang/doan && make build-all"
alias dtest="cd /Users/hant/golang/doan && make test"
```

Sau đó:
```bash
source ~/.zshrc
dev        # Run development
dgen       # Generate code
dbuild     # Build all
dtest      # Run tests
```

## ⚠️ Troubleshooting

### Wire not found
```bash
export PATH=$PATH:$(go env GOPATH)/bin
make install-tools
```

### Import errors
```bash
go clean -modcache
make deps
```

### Wire generation fails
```bash
make clean
make generate
```

### Port in use
Sửa `configs/config.yaml`:
```yaml
http:
  port: 8081  # Change port
```

## 📚 Documentation

- **README.md** - Project overview & full guide
- **DEVELOPMENT.md** - Development guide & best practices
- **RESTRUCTURE_GUIDE.md** - Restructure steps
- **QUICK_REFERENCE.md** - This file
- **Swagger** - http://localhost:8080/api/swagger/index.html

## 🎓 Learning Resources

1. **Wire**: https://github.com/google/wire
2. **Gin**: https://gin-gonic.com/
3. **Go Modules**: https://golang.org/ref/mod
4. **Clean Architecture**: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html

---

**💡 Remember: `make dev` là all you need!**

