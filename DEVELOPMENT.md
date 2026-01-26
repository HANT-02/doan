# Development Guide

## Cấu trúc Monorepo

Project này sử dụng cấu trúc monorepo với nhiều entry points:

- `cmd/http/` - REST API Server
- `cmd/cli/migration/` - Database Migration Tool
- `cmd/cli/get-access-token/` - Token Generator

## Dependency Injection với Wire

### Cách hoạt động

1. **Định nghĩa Providers** - Các hàm khởi tạo dependencies
2. **Wire Configuration** - File `wire.go` định nghĩa cách kết nối dependencies
3. **Generate Code** - Wire tự động tạo code injection trong `wire_gen.go`

### Thêm Dependency Mới

#### Bước 1: Tạo Provider Function

Ví dụ trong `internal/services/product/product_service.go`:

```go
package product

type ProductService struct {
    repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) *ProductService {
    return &ProductService{repo: repo}
}
```

#### Bước 2: Thêm vào Provider Set

Trong `internal/services/provider.go`:

```go
var ProductServiceProvider = wire.NewSet(
    product.NewProductService,
)
```

#### Bước 3: Update Wire Config

Trong `cmd/http/wire.go`:

```go
func wireApp(app *App) error {
    wire.Build(
        database.DBProvider,
        controllers.ControllerProviders,
        usecases.UserUseCaseProviders,
        services.UserServiceProvider,
        services.ProductServiceProvider,  // Thêm provider mới
        inject,
    )
    return nil
}
```

#### Bước 4: Generate Code

```bash
make generate
```

### Best Practices

1. **Một Provider cho mỗi Module** - Tạo file `provider.go` trong mỗi package
2. **Interface over Struct** - Inject interfaces, không phải concrete types
3. **Constructor Pattern** - Sử dụng `New*` functions làm constructors
4. **Provider Sets** - Group related providers bằng `wire.NewSet`

## Layer Architecture

```
Controllers (HTTP Layer)
    ↓
Use Cases (Application Logic)
    ↓
Services (Business Logic)
    ↓
Repositories (Data Access Interface)
    ↓
Infrastructure (Implementation)
```

### Controllers

- Handle HTTP requests/responses
- Validate input
- Call use cases
- Return formatted responses

### Use Cases

- Orchestrate business operations
- Coordinate between services
- Transaction management
- Error handling

### Services

- Core business logic
- Domain rules
- Pure functions (when possible)

### Repositories

- Abstract data access
- Define interfaces only
- No implementation details

### Infrastructure

- Implement repository interfaces
- Database connections
- External service integrations
- Cache implementations

## Testing

### Unit Tests

```go
func TestProductService_Create(t *testing.T) {
    // Arrange
    mockRepo := &MockProductRepository{}
    service := NewProductService(mockRepo)
    
    // Act
    result, err := service.Create(...)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

### Integration Tests

```go
// +build integration

func TestProductRepository_Create(t *testing.T) {
    // Use real database connection
    db := setupTestDB(t)
    defer cleanupTestDB(t, db)
    
    repo := NewProductRepository(db)
    // ... test with real database
}
```

## Configuration Management

### Config File (config.yaml)

```yaml
http:
  host: 0.0.0.0
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: doan
  
redis:
  host: localhost
  port: 6379
```

### Environment Variables

Priority: Environment Variables > Config File > Defaults

```bash
export HTTP_PORT=8080
export DB_HOST=localhost
```

## Database Migrations

### Create Migration

```sql
-- 1_create_product_table.up.sql
CREATE TABLE products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT NOW()
);
```

```sql
-- 1_create_product_table.down.sql
DROP TABLE IF EXISTS products;
```

### Run Migration

```bash
make migrate
```

## API Versioning

### Controller Versioning

```go
// v1/product_controller.go
type ProductControllerV1 struct {}

// v2/product_controller.go
type ProductControllerV2 struct {}
```

### Route Registration

```go
func RegisterRoutesV1(router *gin.Engine, ctrl *ProductControllerV1) {
    v1 := router.Group("/api/v1/products")
    v1.GET("/:id", ctrl.GetByID)
}

func RegisterRoutesV2(router *gin.Engine, ctrl *ProductControllerV2) {
    v2 := router.Group("/api/v2/products")
    v2.GET("/:id", ctrl.GetByID)
}
```

## Error Handling

### Custom Error Types

```go
type AppError struct {
    Code    string
    Message string
    Status  int
}

func (e *AppError) Error() string {
    return e.Message
}
```

### Error Codes

Define in `pkg/error/error_codes.go`:

```go
const (
    ErrNotFound = "NOT_FOUND"
    ErrInvalidInput = "INVALID_INPUT"
    ErrUnauthorized = "UNAUTHORIZED"
)
```

## Logging

### Structured Logging

```go
log.Info("User created", 
    "user_id", userID,
    "email", email,
)

log.Error("Failed to create user",
    "error", err,
    "user_id", userID,
)
```

## Performance Tips

### Database

1. Use connection pooling
2. Add proper indexes
3. Use prepared statements
4. Implement query timeout

### Caching

1. Cache frequently accessed data
2. Set appropriate TTL
3. Implement cache invalidation strategy

### API

1. Implement rate limiting
2. Use compression middleware
3. Optimize JSON serialization
4. Add response caching headers

## Security

### Best Practices

1. **Never commit secrets** - Use environment variables
2. **Validate all inputs** - Sanitize user data
3. **Use HTTPS** - In production
4. **Implement rate limiting** - Prevent abuse
5. **SQL Injection Prevention** - Use parameterized queries
6. **XSS Prevention** - Escape output
7. **CSRF Protection** - Use tokens

## Deployment

### Build for Production

```bash
# Build all binaries
make build-all

# Binaries in bin/ folder
./bin/http-server
./bin/migration
./bin/get-access-token
```

### Docker

```bash
# Build image
make docker-build

# Run container
docker run -p 8080:8080 \
  -e DB_HOST=postgres \
  -e REDIS_HOST=redis \
  doan-app
```

### Kubernetes

```bash
# Deploy with Helm
helm install doan-app deploy/helm/ \
  --set image.tag=v1.0.0 \
  --set database.host=postgres-service
```

## Troubleshooting

### Common Issues

1. **Wire generation fails**
   ```bash
   make clean
   make generate
   ```

2. **Import cycle detected**
   - Review your dependencies
   - Break circular dependencies
   - Use interfaces

3. **Database connection fails**
   - Check config.yaml
   - Verify database is running
   - Check firewall rules

4. **Port already in use**
   - Change port in config
   - Kill existing process

## Resources

- [Wire Documentation](https://github.com/google/wire)
- [Gin Framework](https://gin-gonic.com/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Best Practices](https://golang.org/doc/effective_go)

