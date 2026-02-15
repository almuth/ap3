# Quick Start Guide: AhadPOS 3 Go Implementation

## Overview

This guide provides step-by-step instructions to quickly set up and start implementing the Go-based REST API for AhadPOS 3. Follow these steps to get a working development environment up and running in under 30 minutes.

## Prerequisites

### Required Software
- Go 1.21 or higher
- Git
- SQLite 3
- VS Code (recommended) or your preferred IDE
- Postman or similar API testing tool
- MySQL client (for data export)

### Verify Installation
```bash
# Check Go version
go version

# Check SQLite
sqlite3 --version

# Check Git
git --version
```

## Step 1: Initialize Project (5 minutes)

### 1.1 Create Project Directory
```bash
cd ~/projects  # or your preferred location
mkdir ahadpos-go
cd ahadpos-go
```

### 1.2 Initialize Go Module
```bash
go mod init github.com/ahadmart/ahadpos-go
```

### 1.3 Create Directory Structure
```bash
# Create all directories
mkdir -p cmd/api
mkdir -p internal/{config,controllers,middleware,models,repository,services,utils}
mkdir -p pkg/{database,errors,jwt}
mkdir -p migrations
mkdir -p configs
mkdir -p tests
mkdir -p docs

# Verify structure
tree -L 3
```

## Step 2: Install Dependencies (3 minutes)

### 2.1 Core Dependencies
```bash
# Router and middleware
go get github.com/go-chi/chi/v5
go get github.com/go-chi/cors

# Database
go get gorm.io/gorm
go get gorm.io/driver/sqlite

# Authentication
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt

# Validation
go get github.com/go-playground/validator/v10

# Logging
go get github.com/rs/zerolog

# Configuration
go get github.com/spf13/viper

# Testing
go get github.com/stretchr/testify
```

### 2.2 Development Tools
```bash
# Go tooling
go install github.com/cosmtrek/air@latest  # Hot reload
go install github.com/swaggo/swag/cmd/swag@latest  # Swagger docs

# Verify installations
air -v
swag -v
```

## Step 3: Setup Database Layer (5 minutes)

### 3.1 Create Database Connection
Create file: `pkg/database/database.go`

```go
package database

import (
    "fmt"
    "time"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(dbPath string) error {
    var err error

    DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })

    if err != nil {
        return fmt.Errorf("failed to connect database: %w", err)
    }

    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }

    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)

    // Enable foreign keys
    DB.Exec("PRAGMA foreign_keys = ON")

    return nil
}

func Close() error {
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}
```

### 3.2 Create Initial Schema
Create file: `migrations/000001_init_schema.sql`

```sql
-- Users
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    email TEXT,
    role TEXT NOT NULL DEFAULT 'user',
    status INTEGER NOT NULL DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- Products
CREATE TABLE IF NOT EXISTS barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    barcode TEXT NOT NULL UNIQUE COLLATE NOCASE,
    nama TEXT NOT NULL,
    satuan_id INTEGER NOT NULL,
    status INTEGER NOT NULL DEFAULT 1,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_by INTEGER
);

-- Categories
CREATE TABLE IF NOT EXISTS kategori_barang (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nama TEXT NOT NULL UNIQUE COLLATE NOCASE,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- Sales
CREATE TABLE IF NOT EXISTS penjualan (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    nomor TEXT NOT NULL UNIQUE,
    tanggal TEXT NOT NULL,
    profil_id INTEGER NOT NULL,
    status INTEGER NOT NULL DEFAULT 0,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- Sales Details
CREATE TABLE IF NOT EXISTS penjualan_detail (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    penjualan_id INTEGER NOT NULL,
    barang_id INTEGER NOT NULL,
    qty INTEGER NOT NULL,
    harga_satuan INTEGER NOT NULL,
    total INTEGER NOT NULL,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (penjualan_id) REFERENCES penjualan(id) ON DELETE CASCADE,
    FOREIGN KEY (barang_id) REFERENCES barang(id)
);
```

### 3.3 Run Migration Script
```bash
# Create database and apply schema
sqlite3 ahadpos3.db < migrations/000001_init_schema.sql

# Verify tables created
sqlite3 ahadpos3.db ".tables"
```

## Step 4: Create Configuration (2 minutes)

### 4.1 Config File
Create file: `configs/config.yaml`

```yaml
server:
  port: 8080
  mode: debug  # debug, release
  read_timeout: 30
  write_timeout: 30

database:
  path: ahadpos3.db

jwt:
  secret: your-secret-key-change-this-in-production
  expire_time: 15  # minutes

cors:
  allowed_origins:
    - "*"
  allowed_methods:
    - GET
    - POST
    - PUT
    - DELETE
    - OPTIONS
```

### 4.2 Config Package
Create file: `internal/config/config.go`

```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    CORS     CORSConfig     `mapstructure:"cors"`
}

type ServerConfig struct {
    Port         int    `mapstructure:"port"`
    Mode         string `mapstructure:"mode"`
    ReadTimeout  int    `mapstructure:"read_timeout"`
    WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Path string `mapstructure:"path"`
}

type JWTConfig struct {
    Secret     string `mapstructure:"secret"`
    ExpireTime int    `mapstructure:"expire_time"`
}

type CORSConfig struct {
    AllowedOrigins []string `mapstructure:"allowed_origins"`
    AllowedMethods []string `mapstructure:"allowed_methods"`
}

var AppConfig *Config

func Load(configPath string) error {
    viper.SetConfigFile(configPath)
    viper.SetConfigType("yaml")

    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.mode", "debug")
    viper.SetDefault("database.path", "ahadpos3.db")
    viper.SetDefault("jwt.expire_time", 15)

    if err := viper.ReadInConfig(); err != nil {
        return err
    }

    AppConfig = &Config{}
    return viper.Unmarshal(AppConfig)
}
```

## Step 5: Create Models (3 minutes)

### 5.1 Base Model
Create file: `internal/models/base.go`

```go
package models

import "time"

type BaseModel struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
```

### 5.2 User Model
Create file: `internal/models/user.go`

```go
package models

import "golang.org/x/crypto/bcrypt"

type User struct {
    BaseModel
    Name     string `gorm:"size:100;not null" json:"name" validate:"required"`
    Username string `gorm:"uniqueIndex;size:50;not null" json:"username" validate:"required,max=50"`
    Password string `gorm:"size:255;not null" json:"-" validate:"required,min=6"`
    Email    string `gorm:"size:100" json:"email" validate:"omitempty,email"`
    Role     string `gorm:"size:50;default:'user'" json:"role"`
    Status   int    `gorm:"default:1" json:"status"`
}

func (u *User) SetPassword(password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    u.Password = string(hashedPassword)
    return nil
}

func (u *User) CheckPassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}
```

### 5.3 Product Model
Create file: `internal/models/barang.go`

```go
package models

type Barang struct {
    BaseModel
    Barcode   string `gorm:"uniqueIndex;size:30;not null" json:"barcode" validate:"required,max=30"`
    Nama      string `gorm:"size:100;not null" json:"nama" validate:"required,max=100"`
    SatuanID  uint   `gorm:"not null" json:"satuan_id" validate:"required"`
    Status    int    `gorm:"default:1" json:"status"`
    UpdatedBy uint   `json:"updated_by,omitempty"`
}

const (
    BarangStatusTidakAktif = 0
    BarangStatusAktif     = 1
)
```

## Step 6: Create Error Handling (2 minutes)

### 6.1 Error Package
Create file: `pkg/errors/errors.go`

```go
package errors

import "net/http"

type AppError struct {
    Code       string       `json:"code"`
    Message    string       `json:"message"`
    StatusCode int          `json:"-"`
    Details    []FieldError `json:"details,omitempty"`
}

type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

func (e *AppError) Error() string {
    return e.Message
}

func NewValidationError(field, message string) *AppError {
    return &AppError{
        Code:       "VALIDATION_ERROR",
        Message:    "Validation failed",
        StatusCode: http.StatusBadRequest,
        Details:    []FieldError{{Field: field, Message: message}},
    }
}

func NewNotFoundError(message string) *AppError {
    return &AppError{
        Code:       "NOT_FOUND",
        Message:    message,
        StatusCode: http.StatusNotFound,
    }
}

func NewAuthenticationError(message string) *AppError {
    return &AppError{
        Code:       "AUTHENTICATION_ERROR",
        Message:    message,
        StatusCode: http.StatusUnauthorized,
    }
}
```

## Step 7: Create Response Utilities (2 minutes)

### 7.1 Response Helper
Create file: `internal/utils/response.go`

```go
package utils

import (
    "ahadpos-go/pkg/errors"
    "encoding/json"
    "net/http"
)

type Response struct {
    Success bool         `json:"success"`
    Data    interface{}  `json:"data,omitempty"`
    Error   *ErrorDetail `json:"error,omitempty"`
    Meta    *Meta        `json:"meta,omitempty"`
}

type ErrorDetail struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

type Meta struct {
    Page  int   `json:"page"`
    Limit int   `json:"limit"`
    Total int64 `json:"total"`
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
    response := Response{Success: true, Data: data}
    sendJSON(w, http.StatusOK, response)
}

func SendCreated(w http.ResponseWriter, data interface{}) {
    response := Response{Success: true, Data: data}
    sendJSON(w, http.StatusCreated, response)
}

func SendError(w http.ResponseWriter, err error) {
    statusCode := http.StatusInternalServerError
    errorDetail := &ErrorDetail{Code: "INTERNAL_ERROR", Message: err.Error()}

    if customErr, ok := err.(*errors.AppError); ok {
        statusCode = customErr.StatusCode
        errorDetail = &ErrorDetail{
            Code:    customErr.Code,
            Message: customErr.Message,
            Details: customErr.Details,
        }
    }

    response := Response{Success: false, Error: errorDetail}
    sendJSON(w, statusCode, response)
}

func sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}
```

## Step 8: Create Main Application (5 minutes)

### 8.1 Main Entry Point
Create file: `cmd/api/main.go`

```go
package main

import (
    "ahadpos-go/internal/config"
    "ahadpos-go/pkg/database"
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "github.com/rs/zerolog"
)

func main() {
    // Load configuration
    if err := config.Load("configs/config.yaml"); err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize logger
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Connect to database
    if err := database.Connect(config.AppConfig.Database.Path); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.Close()

    logger.Info().Msg("Database connected successfully")

    // Create router
    r := chi.NewRouter()

    // Middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Timeout(60 * time.Second))

    // CORS middleware
    r.Use(func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")

            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }

            next.ServeHTTP(w, r)
        })
    })

    // Routes
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        response := map[string]string{
            "message": "AhadPOS 3 API",
            "version": "1.0.0",
        }
        utils.SendSuccess(w, response)
    })

    // Start server
    addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
    srv := &http.Server{
        Addr:         addr,
        Handler:      r,
        ReadTimeout:  time.Duration(config.AppConfig.Server.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(config.AppConfig.Server.WriteTimeout) * time.Second,
    }

    logger.Info().Msgf("Starting server on %s", addr)

    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal().Err(err).Msg("Server failed to start")
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Info().Msg("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal().Err(err).Msg("Server forced to shutdown")
    }

    logger.Info().Msg("Server exited")
}
```

### 8.2 Fix Import
Add missing import:
```go
import (
    "ahadpos-go/internal/config"
    "ahadpos-go/internal/utils"  // Add this
    "ahadpos-go/pkg/database"
    // ... rest of imports
)
```

## Step 9: Run and Test (5 minutes)

### 9.1 Run the Application
```bash
# Run directly
go run cmd/api/main.go

# Or use air for hot reload
air init  # First time only
air
```

### 9.2 Test Health Endpoint
```bash
# Terminal 1: Run the application
go run cmd/api/main.go

# Terminal 2: Test
curl http://localhost:8080/health
# Expected output: OK

curl http://localhost:8080/
# Expected output: {"message":"AhadPOS 3 API","version":"1.0.0"}
```

### 9.3 Create Test User
```bash
# Insert test user into database
sqlite3 ahadpos3.db <<EOF
INSERT INTO user (name, username, password, role, status, created_at, updated_at)
VALUES ('Admin User', 'admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', 'admin', 1, datetime('now'), datetime('now'));
EOF

# Verify user created
sqlite3 ahadpos3.db "SELECT id, username, role FROM user;"
```

Note: The password hash above is for "admin123"

## Step 10: Add First API Endpoint (Optional - 5 minutes)

### 10.1 Create Product Controller
Create file: `internal/controllers/barang_controller.go`

```go
package controllers

import (
    "ahadpos-go/internal/utils"
    "ahadpos-go/pkg/database"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi"
)

type BarangController struct{}

func NewBarangController() *BarangController {
    return &BarangController{}
}

func (c *BarangController) RegisterRoutes(r chi.Router) {
    r.Get("/", c.List)
    r.Post("/", c.Create)
    r.Get("/{id}", c.GetByID)
}

type Barang struct {
    ID       uint   `json:"id"`
    Barcode  string `json:"barcode"`
    Nama     string `json:"nama"`
    SatuanID uint   `json:"satuan_id"`
    Status   int    `json:"status"`
}

func (c *BarangController) List(w http.ResponseWriter, r *http.Request) {
    var barang []Barang

    if err := database.DB.Find(&barang).Error; err != nil {
        utils.SendError(w, err)
        return
    }

    utils.SendSuccess(w, barang)
}

func (c *BarangController) GetByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
    if err != nil {
        utils.SendError(w, err)
        return
    }

    var barang Barang
    if err := database.DB.First(&barang, id).Error; err != nil {
        utils.SendError(w, err)
        return
    }

    utils.SendSuccess(w, barang)
}

func (c *BarangController) Create(w http.ResponseWriter, r *http.Request) {
    var barang Barang
    if err := json.NewDecoder(r.Body).Decode(&barang); err != nil {
        utils.SendError(w, err)
        return
    }

    if err := database.DB.Create(&barang).Error; err != nil {
        utils.SendError(w, err)
        return
    }

    utils.SendCreated(w, barang)
}
```

### 10.2 Register Routes in Main

Update `cmd/api/main.go`:
```go
import (
    // ... existing imports
    "ahadpos-go/internal/controllers"  // Add this
)

func main() {
    // ... existing code

    // Register routes
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    // Product routes
    r.Route("/api/v1/barang", func(r chi.Router) {
        barangController := controllers.NewBarangController()
        barangController.RegisterRoutes(r)
    })

    // ... rest of code
}
```

### 10.3 Test Product API
```bash
# List products (empty initially)
curl http://localhost:8080/api/v1/barang

# Create a product
curl -X POST http://localhost:8080/api/v1/barang \
  -H "Content-Type: application/json" \
  -d '{
    "barcode": "123456789",
    "nama": "Test Product",
    "satuan_id": 1,
    "status": 1
  }'

# Get product by ID
curl http://localhost:8080/api/v1/barang/1

# List products (should have one now)
curl http://localhost:8080/api/v1/barang
```

## Next Steps

### Immediate Next Actions
1. **Add Authentication**: Implement JWT-based authentication
2. **Add More Controllers**: Sales, purchases, customers, etc.
3. **Add Validation**: Use validator package for input validation
4. **Add Tests**: Write unit and integration tests
5. **Add Documentation**: Set up Swagger for API documentation

### Recommended Development Order
1. ✅ Project setup and basic structure
2. ✅ Database connection and basic models
3. ⏳ Authentication & Authorization
4. ⏳ Product/Inventory management
5. ⏳ Sales management
6. ⏳ Purchase management
7. ⏳ POS/Cashier features
8. ⏳ Reporting
9. ⏳ Testing and optimization

### Learning Resources
- Go Documentation: https://golang.org/doc/
- Chi Router: https://github.com/go-chi/chi
- GORM: https://gorm.io/docs/
- REST API Design: https://restfulapi.net/

## Common Issues and Solutions

### Issue: Import Errors
**Problem**: `cannot find package`
**Solution**: Run `go mod tidy` to download dependencies

### Issue: Database Locked
**Problem**: "database is locked"
**Solution**: Close all connections, check for other processes using the database

### Issue: Port Already in Use
**Problem**: "bind: address already in use"
**Solution**: Change port in configs/config.yaml or kill the process using port 8080

### Issue: CORS Errors
**Problem**: Browser blocks API requests
**Solution**: Ensure CORS middleware is properly configured

## Development Workflow

### Daily Development
```bash
# 1. Pull latest changes
git pull origin main

# 2. Start server with hot reload
air

# 3. Run tests
go test ./...

# 4. Format code
go fmt ./...

# 5. Check for issues
go vet ./...
```

### Commit Code
```bash
# Stage changes
git add .

# Commit with message
git commit -m "feat: add product CRUD operations"

# Push
git push origin feature/your-feature-name
```

## Testing Strategy

### Unit Tests
```go
// tests/barang_test.go
package tests

import (
    "ahadpos-go/internal/models"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestBarangModel(t *testing.T) {
    barang := models.Barang{
        Barcode:  "123456789",
        Nama:     "Test Product",
        SatuanID: 1,
        Status:   1,
    }

    assert.Equal(t, "123456789", barang.Barcode)
    assert.Equal(t, "Test Product", barang.Nama)
    assert.Equal(t, models.BarangStatusAktif, barang.Status)
}
```

### Run Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v -run TestBarangModel ./tests
```

## Deployment Preparation

### Build Binary
```bash
# Build for current OS
go build -o ahadpos-api cmd/api/main.go

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o ahadpos-api-linux cmd/api/main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o ahadpos-api.exe cmd/api/main.go
```

### Docker Build
```bash
# Build Docker image
docker build -t ahadpos-api:latest .

# Run container
docker run -p 8080:8080 -v $(pwd)/data:/app/data ahadpos-api:latest
```

## Quick Reference

### Essential Commands
```bash
# Run application
go run cmd/api/main.go

# Build application
go build -o ahadpos-api cmd/api/main.go

# Run tests
go test ./...

# Format code
go fmt ./...

# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...

# Generate Swagger docs
swag init -g cmd/api/main.go -o docs
```

### Directory Reference
```
ahadpos-go/
├── cmd/api/          # Application entry point
├── internal/          # Private application code
│   ├── config/        # Configuration management
│   ├── controllers/   # HTTP handlers
│   ├── middleware/    # HTTP middleware
│   ├── models/        # Data models
│   ├── repository/    # Data access layer
│   ├── services/      # Business logic
│   └── utils/         # Utility functions
├── pkg/               # Public libraries
│   ├── database/      # Database connection
│   ├── errors/        # Error handling
│   └── jwt/           # JWT utilities
├── migrations/        # Database migrations
├── configs/           # Configuration files
├── tests/             # Test files
└── docs/              # Documentation
```

## Conclusion

Congratulations! You now have a working Go-based REST API foundation for AhadPOS 3. The application is running, the database is connected, and you have tested the first endpoint. From here, follow the development order to implement the remaining features.

For more detailed implementation guidance, refer to:
- `MIGRATION_PLAN.md` - Overall migration strategy
- `TECHNICAL_GUIDE.md` - Detailed implementation patterns
- `DATABASE_MIGRATION_GUIDE.md` - Database migration instructions

Good luck with your implementation! 🚀
