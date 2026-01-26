# Scripts Directory

Collection of utility scripts for the project.

## Available Scripts

### 1. generate.sh
Tự động generate tất cả code (Wire + Swagger)

**Usage:**
```bash
./scripts/generate.sh
```

**What it does:**
- Install Wire and Swag if not found
- Generate Wire dependency injection code for all applications
- Generate Swagger documentation

**Equivalent to:** `make generate`

---

### 2. quick-setup.sh
Setup nhanh cho lần đầu tiên sử dụng project

**Usage:**
```bash
./scripts/quick-setup.sh
```

**What it does:**
- Install required tools (Wire, Swag)
- Download Go dependencies
- Copy config file from sample
- Optionally start Docker services
- Generate all code

**Equivalent to:** Manual setup steps in SETUP.md

---

### 3. check-env.sh
Kiểm tra môi trường development

**Usage:**
```bash
./scripts/check-env.sh
```

**What it does:**
- Check Go installation
- Check Wire installation
- Check Swag installation
- Check config file exists
- Check PostgreSQL
- Check Docker
- Check generated files

---

## Making Scripts Executable

```bash
chmod +x scripts/*.sh
```

## Using Scripts

### First Time Setup
```bash
./scripts/quick-setup.sh
```

### Check Environment
```bash
./scripts/check-env.sh
```

### Generate Code
```bash
./scripts/generate.sh
```

## Or Use Makefile

All scripts can be executed via Makefile:

```bash
make install-tools    # Install Wire & Swag
make generate        # Generate all code
make deps           # Download dependencies
```

## Troubleshooting

### Permission Denied
```bash
chmod +x scripts/script-name.sh
```

### Script Not Found
Make sure you're in the project root directory:
```bash
cd /Users/hant/golang/doan
./scripts/script-name.sh
```

### Command Not Found (wire, swag)
Add Go bin to PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Or add to `~/.zshrc`:
```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

