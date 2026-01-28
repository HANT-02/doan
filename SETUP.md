# Setup Guide - Hướng dẫn thiết lập Project

## 📋 Checklist

- [ ] Cài đặt Go 1.21+
- [ ] Cài đặt Wire
- [ ] Cài đặt Swag
- [ ] Download dependencies
- [ ] Setup database (PostgreSQL)
- [ ] Setup Redis (optional)
- [ ] Copy config file
- [ ] Generate code
- [ ] Run application

## Bước 1: Kiểm tra Go

```bash
go version
# Phải >= 1.21
```

Nếu chưa có Go: https://golang.org/dl/

## Bước 2: Cài đặt Tools

### Cách 1: Tự động (Khuyến nghị)

```bash
make install-tools
```

### Cách 2: Thủ công

```bash
# Wire
go install github.com/google/wire/cmd/wire@latest

# Swag
go install github.com/swaggo/swag/cmd/swag@latest

# Verify
wire --version
swag --version
```

### Thêm Go bin vào PATH (nếu cần)

```bash
# Thêm vào ~/.zshrc hoặc ~/.bashrc
export PATH=$PATH:$(go env GOPATH)/bin

# Reload
source ~/.zshrc
```

## Bước 3: Download Dependencies

```bash
make deps
```

Hoặc:

```bash
go mod download
go mod tidy
```

## Bước 4: Setup Database

### Option 1: Sử dụng Docker (Khuyến nghị)

```bash
# Start PostgreSQL và Redis
make docker-local-up

# Verify
docker ps
```

### Option 2: Cài đặt local

#### PostgreSQL

**macOS:**
```bash
brew install postgresql@15
brew services start postgresql@15
createdb doan
```

**Linux:**
```bash
sudo apt-get install postgresql-15
sudo systemctl start postgresql
sudo -u postgres createdb doan
```

#### Redis (Optional)

**macOS:**
```bash
brew install redis
brew services start redis
```

**Linux:**
```bash
sudo apt-get install redis-server
sudo systemctl start redis
```

## Bước 5: Configuration

### Copy config file

```bash
cp configs/config.yaml.sample configs/config.yaml
```

### Edit config.yaml

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
  sslmode: disable
  
redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
```

### Environment Variables (Optional)

```bash
# Tạo .env file
cat > .env << 'EOF'
HTTP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=doan
REDIS_HOST=localhost
REDIS_PORT=6379
EOF
```

## Bước 6: Generate Code

```bash
make generate
```

Lệnh này sẽ:
1. Generate Wire dependency injection code
2. Generate Swagger documentation

## Bước 7: Run Migrations

```bash
make migrate
```

Hoặc:

```bash
go run cmd/cli/migration/main.go cmd/cli/migration/wire_gen.go
```

## Bước 8: Run Application

```bash
make dev
```

Hoặc:

```bash
make run
```

## Bước 9: Verify

### Check HTTP server

```bash
curl http://localhost:8080/ping
# Response: {"message":"pong"}
```

### Check Swagger docs

Mở browser: http://localhost:8080/api/swagger/index.html

## 🎉 Xong! Application đã chạy

---

## Troubleshooting

### 1. Wire not found

```bash
# Check PATH
echo $PATH | grep $(go env GOPATH)/bin

# If not found, add to PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Add to ~/.zshrc for permanent
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

### 2. Cannot connect to database

```bash
# Check PostgreSQL is running
psql -U postgres -c "SELECT 1"

# Check connection
psql -U postgres -d doan -c "SELECT 1"

# Check config.yaml
cat configs/config.yaml | grep -A5 database
```

### 3. Port already in use

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in config.yaml
```

### 4. Import errors

```bash
# Clean mod cache
go clean -modcache

# Re-download
make deps

# Or
go mod download
go mod tidy
```

### 5. Wire generation fails

```bash
# Clean and regenerate
make clean
make generate

# Manual check
cd cmd/http
wire
cd ../..
```

### 6. Permission denied on scripts

```bash
chmod +x scripts/*.sh
```

---

## Next Steps

### 1. Development

```bash
# Start development
make dev
```

### 2. Read Documentation

- `README.md` - Project overview
- `QUICK_REFERENCE.md` - Quick commands
- `DEVELOPMENT.md` - Development guide

### 3. Explore Code

```
cmd/http/               # HTTP server
internal/entities/      # Domain models
internal/services/      # Business logic
internal/usecases/      # Application logic
```

### 4. Create Feature

Follow workflow in `DEVELOPMENT.md`

---

## Useful Commands

```bash
# Development
make dev              # Auto-generate and run
make run              # Generate and run
make run-no-gen       # Run without regenerate

# Code Generation
make generate         # Generate all
make wire            # Generate Wire only
make swagger         # Generate Swagger only

# Build
make build           # Build HTTP server
make build-all       # Build all binaries

# Testing
make test            # Run tests
make test-coverage   # Test with coverage

# Database
make migrate         # Run migrations

# Docker
make docker-local-up    # Start services
make docker-local-down  # Stop services

# Clean
make clean           # Clean generated files

# Help
make help            # Show all commands
```

---

## Development Workflow

```bash
# Morning routine
cd /path/to/doan
make docker-local-up
make dev

# After making changes
make generate
make run

# Before commit
make test
make fmt

# End of day
make docker-local-down
```

---

## Quick Start (TL;DR)

```bash
# Install tools
make install-tools

# Setup
make deps
make docker-local-up
cp configs/config.yaml.sample configs/config.yaml

# Run
make dev
```

---

## Support

If you encounter any issues:

1. Check `TROUBLESHOOTING` section above
2. Check `README.md` for detailed documentation
3. Check `DEVELOPMENT.md` for development guide
4. Check existing issues in Git repository

---

**Happy Coding! 🚀**

