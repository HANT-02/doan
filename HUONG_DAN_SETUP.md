# Hướng Dẫn Setup và Sử Dụng Authentication Flow

## Tổng Quan
Dự án này đã được cấu trúc lại theo Clean Architecture với:
- **GORM** cho database ORM
- **Wire** cho dependency injection
- **JWT** cho authentication
- **Swagger** cho API documentation

## Các Bước Setup

### 1. Cài Đặt Dependencies

```bash
# Cài đặt công cụ cần thiết
make install-tools

# Tải dependencies
go mod tidy
```

### 2. Cấu Hình Database

Chỉnh sửa `configs/config.yaml`:

```yaml
database:
  driver: postgres
  user: root          # Đổi thành user của bạn
  password: password  # Đổi thành password của bạn
  protocol: tcp
  host: localhost
  port: 5432
  schema: doan

jwt:
  secret: "your-secret-key-change-this-in-production"
  access_token_duration: 24h
  refresh_token_duration: 168h
```

### 3. Tạo Database

```bash
# Tạo database PostgreSQL
createdb doan

# Hoặc dùng psql
psql -U postgres -c "CREATE DATABASE doan;"
```

### 4. Chạy Migration

```bash
# Generate Wire code và chạy migration
make wire && make migrate
```

Migration sẽ tự động tạo các bảng cần thiết bằng GORM AutoMigrate.

### 5. Seed Dữ Liệu Mẫu

```bash
# Tạo user mẫu để test
make seed
```

Lệnh này sẽ tạo 5 user mẫu:
- **admin@example.com** (ADMIN)
- **teacher@example.com** (TEACHER)
- **student1@example.com** (STUDENT - Active)
- **student2@example.com** (STUDENT - Active)
- **student3@example.com** (STUDENT - Inactive)

**Password mặc định cho tất cả:** `password123`

### 6. Chạy Server

```bash
# Chạy HTTP server
make dev
```

Server sẽ chạy tại: `http://localhost:9000`

### 7. Xem Swagger Documentation

Mở trình duyệt và truy cập:
```
http://localhost:9000/api/swagger/index.html
```

## API Endpoints

### Authentication v1

**Base URL:** `/v1/auth`

#### 1. Login
```
POST /v1/auth/login
Content-Type: application/json

{
  "username": "admin@example.com",
  "password": "password123"
}
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOi...",
    "refresh_token": "eyJhbGciOi...",
    "user": {
      "id": "uuid",
      "code": "ADMIN001",
      "full_name": "Admin User",
      "email": "admin@example.com",
      "role": "ADMIN",
      "is_active": true
    }
  }
}
```

#### 2. Logout
```
POST /v1/auth/logout
Content-Type: application/json

{
  "token": "your_access_token"
}
```

#### 3. Refresh Token
```
POST /v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "your_refresh_token"
}
```

### Authentication v2

Tương tự v1 nhưng base URL là `/v2/auth`

## Test API với cURL

### Login
```bash
curl -X POST http://localhost:9000/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin@example.com",
    "password": "password123"
  }'
```

### Sử dụng Token để truy cập Protected Route
```bash
curl -X GET http://localhost:9000/api/protected-endpoint \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### Refresh Token
```bash
curl -X POST http://localhost:9000/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

## Cấu Trúc Dự Án

```
doan/
├── cmd/
│   ├── http/              # HTTP server
│   │   ├── main.go
│   │   ├── wire.go        # Wire DI config
│   │   ├── controllers/   # HTTP handlers
│   │   └── middleware/    # Middleware (auth, cors, logger)
│   └── cli/
│       ├── migration/     # Database migration
│       └── seed/          # Data seeding
├── internal/
│   ├── entities/          # Domain models (GORM)
│   ├── repositories/      # Data access layer
│   ├── services/          # Business services
│   ├── usecases/          # Use cases
│   └── infrastructure/    # External dependencies
├── pkg/                   # Shared packages
│   ├── utils/            # Utilities (JWT, etc.)
│   ├── config/           # Configuration
│   ├── logger/           # Logging
│   └── types/            # Common types
└── configs/
    └── config.yaml       # App configuration
```

## Luồng Authentication

### Login Flow
1. Client gửi username/password
2. Controller nhận request → gọi LoginUseCase
3. LoginUseCase → gọi AuthService
4. AuthService:
   - Tìm user trong DB (qua UserRepository)
   - Verify password (bcrypt)
   - Generate JWT tokens (access + refresh)
5. Trả về tokens và thông tin user

### Protected Routes
Để bảo vệ routes với JWT:

```go
import "doan/cmd/http/middleware"

protected := router.Group("/api")
protected.Use(middleware.AuthMiddleware(configManager))
{
    protected.GET("/profile", controller.GetProfile)
}

// Với role-based access
admin := router.Group("/admin")
admin.Use(
    middleware.AuthMiddleware(configManager),
    middleware.RoleMiddleware("ADMIN"),
)
{
    admin.GET("/users", controller.ListUsers)
}
```

## Wire Dependency Injection

### Khi nào cần regenerate Wire?

- Thêm dependency mới vào constructor
- Thay đổi provider
- Thêm repository/service/usecase mới

```bash
# Regenerate wire code
make wire
```

### Provider Files

- `internal/infrastructure/database/provider.go` - Database providers
- `internal/services/provider.go` - Service providers
- `internal/usecases/provider.go` - UseCase providers
- `cmd/http/controllers/provider.go` - Controller providers

## Makefile Commands

```bash
make install-tools  # Cài đặt Wire, Swag, GORM
make wire          # Generate Wire DI code
make swagger       # Generate Swagger docs
make generate      # Generate all (wire + swagger)
make migrate       # Chạy database migration
make seed          # Seed dữ liệu mẫu
make dev           # Chạy server (development)
make run           # Chạy server (production)
make build         # Build binary
make test          # Chạy tests
make clean         # Xóa generated files
```

## Thêm Entity Mới

### 1. Tạo Entity
`internal/entities/your_entity.go`
```go
type YourEntity struct {
    ID        string         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
    Name      string         `gorm:"type:varchar(255)"`
    CreatedAt time.Time      `gorm:"default:now()"`
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### 2. Tạo Repository Interface
`internal/repositories/interface/your_entity.go`
```go
type YourEntityRepository interface {
    repositories.BaseRepository[entities.YourEntity]
}
```

### 3. Implement Repository
`internal/infrastructure/database/postgres/implement/your_entity_repository.go`
```go
func NewYourEntityRepository(
    db *gorm.DB,
    log logger.Logger,
    manager config.Manager,
) repointerface.YourEntityRepository {
    modelRepo := postgres.NewBaseRepository[entities.YourEntity](log, manager, db, "your_entities")
    return &yourEntityRepository{
        BaseDependency: base_struct.BaseDependency{
            Log:           log,
            ConfigManager: manager,
        },
        BaseRepository: modelRepo,
        db:             db,
    }
}
```

### 4. Thêm vào Provider
`internal/infrastructure/database/provider.go`
```go
var DBProvider = wire.NewSet(
    GetDBContext,
    implement.NewUserRepository,
    implement.NewYourEntityRepository, // Thêm dòng này
)
```

### 5. Thêm vào Migration
`internal/infrastructure/database/postgres/migration.go`
```go
entities := []interface{}{
    &entities.User{},
    &entities.YourEntity{}, // Thêm dòng này
}
```

### 6. Generate Wire và Migrate
```bash
make wire && make migrate
```

## Troubleshooting

### Lỗi Wire Generation
```bash
# Xóa generated files và regenerate
rm -f cmd/http/wire_gen.go cmd/cli/migration/wire_gen.go
make wire
```

### Lỗi Migration
```bash
# Enable UUID extension thủ công
psql -U root -d doan -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"

# Retry migration
make migrate
```

### Lỗi JWT Token Invalid
- Kiểm tra JWT secret trong config.yaml
- Verify token chưa hết hạn
- Đảm bảo format token: `Bearer <token>`

## Security Best Practices

1. **Đổi JWT Secret** trong production
2. **Sử dụng HTTPS** trong production
3. **Implement Rate Limiting** cho auth endpoints
4. **Add Token Blacklist** với Redis cho logout
5. **Validate tất cả inputs**
6. **Hash passwords** với bcrypt (đã implement)

## Next Steps

- [x] Hoàn thiện authentication flow (login/logout/refresh)
- [x] Tạo middleware authentication
- [x] Generate Swagger documentation
- [x] Setup Wire dependency injection
- [ ] Implement token blacklist với Redis
- [ ] Add email verification
- [ ] Add password reset
- [ ] Implement 2FA
- [ ] Add OAuth (Google, Facebook)

## Tài Liệu Chi Tiết

Xem `AUTH_FLOW.md` để biết thêm chi tiết về:
- Architecture design
- Security implementation
- Advanced usage
- Contributing guidelines

## Support

Nếu gặp vấn đề, kiểm tra:
1. Database connection trong config.yaml
2. PostgreSQL đang chạy
3. Go version >= 1.25
4. Đã chạy `make install-tools`
5. Đã chạy `make wire` và `make migrate`

