# AhadPOS 3 Go Implementation - Technical Guide

## 1. Project Initialization

### Initialize Go Module
```bash
mkdir ahadpos-go
cd ahadpos-go
go mod init github.com/ahadmart/ahadpos-go
```

### Install Dependencies
```bash
# Core dependencies
go get github.com/go-chi/chi/v5
go get github.com/go-chi/cors
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/golang-jwt/jwt/v5
go get github.com/go-playground/validator/v10
go get github.com/rs/zerolog
go get github.com/spf13/viper

# Testing
go get github.com/stretchr/testify
go get github.com/golang/mock/mockgen

# Documentation
go get github.com/swaggo/swag/cmd/swag
go get github.com/swaggo/files
go get github.com/swaggo/gin-swagger
```

## 2. Database Layer Setup

### Database Connection (pkg/database/database.go)
```go
package database

import (
    "database/sql"
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

    // Enable foreign keys
    sqlDB, err := DB.DB()
    if err != nil {
        return fmt.Errorf("failed to get database instance: %w", err)
    }

    // Set connection pool settings
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

### Model Base Structure (internal/models/base.go)
```go
package models

import (
    "time"
)

type BaseModel struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type TimestampModel struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
    UpdatedBy uint      `json:"updated_by"`
}
```

## 3. Core Models

### Product Model (internal/models/barang.go)
```go
package models

type Barang struct {
    ID                uint   `gorm:"primarykey" json:"id"`
    Barcode           string `gorm:"uniqueIndex;size:30;not null" json:"barcode" validate:"required,max=30"`
    Nama              string `gorm:"size:100;not null" json:"nama" validate:"required,max=100"`
    StrukturID        *uint  `gorm:"index" json:"struktur_id,omitempty"`
    KategoriID        *uint  `gorm:"index" json:"kategori_id,omitempty"`
    SatuanID          uint   `gorm:"not null;index" json:"satuan_id" validate:"required"`
    RakID             *uint  `gorm:"index" json:"rak_id,omitempty"`
    RestockPoint      int    `gorm:"default:0" json:"restock_point"`
    RestockLevel      int    `gorm:"default:0" json:"restock_level"`
    RestockMin        int    `gorm:"default:0" json:"restock_min"`
    VariantCoefficient float64 `json:"variant_coefficient"`
    Status            int    `gorm:"default:1;comment:0=tidak aktif,1=aktif" json:"status"`
    UpdatedBy         uint   `json:"updated_by"`

    // Relations
    Kategori         *KategoriBarang `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
    Satuan           *SatuanBarang   `gorm:"foreignKey:SatuanID" json:"satuan,omitempty"`
    Rak              *RakBarang      `gorm:"foreignKey:RakID" json:"rak,omitempty"`
    Struktur         *StrukturBarang `gorm:"foreignKey:StrukturID" json:"struktur,omitempty"`
    HargaJual        []HargaJual     `gorm:"foreignKey:BarangID" json:"harga_jual,omitempty"`
    InventoryBalance *InventoryBalance `gorm:"foreignKey:BarangID" json:"inventory_balance,omitempty"`
}

const (
    BarangStatusTidakAktif = 0
    BarangStatusAktif     = 1
)
```

### User Model (internal/models/user.go)
```go
package models

import "golang.org/x/crypto/bcrypt"

type User struct {
    ID        uint   `gorm:"primarykey" json:"id"`
    Name      string `gorm:"size:100;not null" json:"name" validate:"required"`
    Username  string `gorm:"uniqueIndex;size:50;not null" json:"username" validate:"required,max=50"`
    Password  string `gorm:"size:255;not null" json:"-" validate:"required,min=6"`
    Email     string `gorm:"size:100" json:"email" validate:"omitempty,email"`
    Role      string `gorm:"size:50;default:'user'" json:"role"`
    Status    int    `gorm:"default:1" json:"status"`
    MenuID    *uint  `json:"menu_id"`
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

### Sale Model (internal/models/penjualan.go)
```go
package models

type Penjualan struct {
    BaseModel
    Nomor          string  `gorm:"size:45;uniqueIndex" json:"nomor" validate:"required,max=45"`
    Tanggal        string  `gorm:"size:19;not null" json:"tanggal" validate:"required"`
    ProfilID       uint    `gorm:"not null;index" json:"profil_id" validate:"required"`
    HutangPiutangID *uint   `gorm:"index" json:"hutang_piutang_id,omitempty"`
    TransferMode   int     `gorm:"default:0" json:"transfer_mode"`
    Status         int     `gorm:"default:0" json:"status"`
    UpdatedBy      uint    `json:"updated_by"`

    // Relations
    Profil           *Profil          `gorm:"foreignKey:ProfilID" json:"profil,omitempty"`
    HutangPiutang    *HutangPiutang   `gorm:"foreignKey:HutangPiutangID" json:"hutang_piutang,omitempty"`
    Details          []PenjualanDetail `gorm:"foreignKey:PenjualanID" json:"details,omitempty"`
    Diskon           []PenjualanDiskon `gorm:"foreignKey:PenjualanID" json:"diskon,omitempty"`
}

const (
    PenjualanStatusDraft   = 0
    PenjualanStatusPiutang = 1
    PenjualanStatusLunas   = 2
)

type PenjualanDetail struct {
    BaseModel
    PenjualanID  uint    `gorm:"not null;index" json:"penjualan_id" validate:"required"`
    BarangID     uint    `gorm:"not null;index" json:"barang_id" validate:"required"`
    Qty          int     `gorm:"not null" json:"qty" validate:"required,gt=0"`
    HargaSatuan  int64   `gorm:"not null" json:"harga_satuan" validate:"required,gt=0"` // In cents
    Diskon       int64   `gorm:"default:0" json:"diskon"` // In cents
    DiskonPersen float64 `gorm:"default:0" json:"diskon_persen"`
    Total        int64   `gorm:"not null" json:"total" validate:"required,gt=0"` // In cents
    UpdatedBy    uint    `json:"updated_by"`

    // Relations
    Barang *Barang `gorm:"foreignKey:BarangID" json:"barang,omitempty"`
}
```

## 4. Repository Pattern

### Base Repository (internal/repository/base.go)
```go
package repository

import (
    "ahadpos-go/internal/models"
    "ahadpos-go/pkg/database"
    "gorm.io/gorm"
    "gorm.io/gorm/clause"
)

type BaseRepository struct {
    DB *gorm.DB
}

func NewBaseRepository() *BaseRepository {
    return &BaseRepository{DB: database.DB}
}

func (r *BaseRepository) GetDB() *gorm.DB {
    return r.DB
}

func (r *BaseRepository) Create(model interface{}) error {
    return r.DB.Create(model).Error
}

func (r *BaseRepository) Update(model interface{}) error {
    return r.DB.Save(model).Error
}

func (r *BaseRepository) Delete(model interface{}) error {
    return r.DB.Delete(model).Error
}

func (r *BaseRepository) First(model interface{}, conditions ...interface{}) error {
    return r.DB.First(model, conditions...).Error
}

func (r *BaseRepository) Find(model interface{}, conditions ...interface{}) error {
    return r.DB.Find(model, conditions...).Error
}

func (r *BaseRepository) WithPreloads(preloads []string) *gorm.DB {
    db := r.DB
    for _, preload := range preloads {
        db = db.Preload(preload)
    }
    return db
}
```

### Product Repository (internal/repository/barang_repository.go)
```go
package repository

import (
    "ahadpos-go/internal/models"
)

type BarangRepository struct {
    *BaseRepository
}

func NewBarangRepository() *BarangRepository {
    return &BarangRepository{
        BaseRepository: NewBaseRepository(),
    }
}

func (r *BarangRepository) Create(barang *models.Barang) error {
    return r.DB.Create(barang).Error
}

func (r *BarangRepository) GetByID(id uint, preloads []string) (*models.Barang, error) {
    var barang models.Barang
    db := r.DB

    for _, preload := range preloads {
        db = db.Preload(preload)
    }

    err := db.First(&barang, id).Error
    if err != nil {
        return nil, err
    }
    return &barang, nil
}

func (r *BarangRepository) GetByBarcode(barcode string) (*models.Barang, error) {
    var barang models.Barang
    err := r.DB.Where("barcode = ?", barcode).First(&barang).Error
    if err != nil {
        return nil, err
    }
    return &barang, nil
}

func (r *BarangRepository) List(page, limit int, search string, status *int) ([]models.Barang, int64, error) {
    var barang []models.Barang
    var total int64

    db := r.DB.Model(&models.Barang{})

    // Apply filters
    if search != "" {
        db = db.Where("barcode LIKE ? OR nama LIKE ?", "%"+search+"%", "%"+search+"%")
    }

    if status != nil {
        db = db.Where("status = ?", *status)
    }

    // Count total
    if err := db.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Apply pagination
    offset := (page - 1) * limit
    if err := db.Offset(offset).Limit(limit).Find(&barang).Error; err != nil {
        return nil, 0, err
    }

    return barang, total, nil
}

func (r *BarangRepository) Update(barang *models.Barang) error {
    return r.DB.Save(barang).Error
}

func (r *BarangRepository) Delete(id uint) error {
    return r.DB.Delete(&models.Barang{}, id).Error
}

func (r *BarangRepository) UpdateStock(id uint, qtyChange int) error {
    return r.DB.Model(&models.InventoryBalance{}).
        Where("barang_id = ?", id).
        UpdateColumn("qty", gorm.Expr("qty + ?", qtyChange)).Error
}
```

## 5. Service Layer

### Product Service (internal/services/barang_service.go)
```go
package services

import (
    "ahadpos-go/internal/models"
    "ahadpos-go/internal/repository"
    "ahadpos-go/pkg/errors"
    "fmt"
)

type BarangService struct {
    repo *repository.BarangRepository
}

func NewBarangService() *BarangService {
    return &BarangService{
        repo: repository.NewBarangRepository(),
    }
}

func (s *BarangService) Create(barang *models.Barang) error {
    // Check if barcode exists
    existing, err := s.repo.GetByBarcode(barang.Barcode)
    if err == nil && existing != nil {
        return errors.NewConflictError("barcode already exists")
    }

    // Validate barang
    if err := validateBarang(barang); err != nil {
        return err
    }

    return s.repo.Create(barang)
}

func (s *BarangService) GetByID(id uint) (*models.Barang, error) {
    preloads := []string{"Kategori", "Satuan", "Rak", "HargaJual"}
    barang, err := s.repo.GetByID(id, preloads)
    if err != nil {
        return nil, errors.NewNotFoundError("barang not found")
    }
    return barang, nil
}

func (s *BarangService) GetByBarcode(barcode string) (*models.Barang, error) {
    barang, err := s.repo.GetByBarcode(barcode)
    if err != nil {
        return nil, errors.NewNotFoundError("barang not found")
    }
    return barang, nil
}

func (s *BarangService) List(page, limit int, search string, status *int) ([]models.Barang, int64, error) {
    return s.repo.List(page, limit, search, status)
}

func (s *BarangService) Update(id uint, barang *models.Barang) error {
    // Check if exists
    existing, err := s.repo.GetByID(id, nil)
    if err != nil {
        return errors.NewNotFoundError("barang not found")
    }

    // Check if barcode is taken by another record
    if barang.Barcode != existing.Barcode {
        other, err := s.repo.GetByBarcode(barang.Barcode)
        if err == nil && other != nil && other.ID != id {
            return errors.NewConflictError("barcode already exists")
        }
    }

    barang.ID = id
    return s.repo.Update(barang)
}

func (s *BarangService) Delete(id uint) error {
    _, err := s.repo.GetByID(id, nil)
    if err != nil {
        return errors.NewNotFoundError("barang not found")
    }
    return s.repo.Delete(id)
}

func validateBarang(barang *models.Barang) error {
    if barang.Barcode == "" {
        return errors.NewValidationError("barcode", "barcode is required")
    }
    if barang.Nama == "" {
        return errors.NewValidationError("nama", "nama is required")
    }
    if barang.SatuanID == 0 {
        return errors.NewValidationError("satuan_id", "satuan_id is required")
    }
    return nil
}
```

### Authentication Service (internal/services/auth_service.go)
```go
package services

import (
    "ahadpos-go/internal/models"
    "ahadpos-go/internal/repository"
    "ahadpos-go/pkg/errors"
    "ahadpos-go/pkg/jwt"
    "time"
)

type AuthService struct {
    userRepo *repository.UserRepository
}

func NewAuthService() *AuthService {
    return &AuthService{
        userRepo: repository.NewUserRepository(),
    }
}

type LoginRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
    Token        string `json:"token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"`
    User         *models.User `json:"user"`
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
    // Find user by username
    user, err := s.userRepo.GetByUsername(req.Username)
    if err != nil {
        return nil, errors.NewAuthenticationError("invalid credentials")
    }

    // Check password
    if !user.CheckPassword(req.Password) {
        return nil, errors.NewAuthenticationError("invalid credentials")
    }

    // Check status
    if user.Status != 1 {
        return nil, errors.NewAuthenticationError("user is inactive")
    }

    // Generate tokens
    token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
    if err != nil {
        return nil, err
    }

    refreshToken, err := jwt.GenerateRefreshToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &LoginResponse{
        Token:        token,
        RefreshToken: refreshToken,
        ExpiresIn:    15 * 60, // 15 minutes
        User:         user,
    }, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*LoginResponse, error) {
    claims, err := jwt.ValidateRefreshToken(refreshToken)
    if err != nil {
        return nil, errors.NewAuthenticationError("invalid refresh token")
    }

    user, err := s.userRepo.GetByID(claims.UserID)
    if err != nil {
        return nil, errors.NewAuthenticationError("user not found")
    }

    // Generate new tokens
    token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
    if err != nil {
        return nil, err
    }

    newRefreshToken, err := jwt.GenerateRefreshToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &LoginResponse{
        Token:        token,
        RefreshToken: newRefreshToken,
        ExpiresIn:    15 * 60,
        User:         user,
    }, nil
}

func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return nil, errors.NewNotFoundError("user not found")
    }
    return user, nil
}
```

## 6. HTTP Handlers

### Response Helper (internal/utils/response.go)
```go
package utils

import (
    "ahadpos-go/pkg/errors"
    "net/http"
)

type Response struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data,omitempty"`
    Error   *ErrorDetail `json:"error,omitempty"`
    Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorDetail struct {
    Code    string              `json:"code"`
    Message string              `json:"message"`
    Details []FieldError        `json:"details,omitempty"`
}

type FieldError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type Meta struct {
    Page  int   `json:"page"`
    Limit int   `json:"limit"`
    Total int64 `json:"total"`
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
    response := Response{
        Success: true,
        Data:    data,
    }
    SendJSON(w, http.StatusOK, response)
}

func SendCreated(w http.ResponseWriter, data interface{}) {
    response := Response{
        Success: true,
        Data:    data,
    }
    SendJSON(w, http.StatusCreated, response)
}

func SendPaginated(w http.ResponseWriter, data interface{}, page, limit int, total int64) {
    response := Response{
        Success: true,
        Data:    data,
        Meta: &Meta{
            Page:  page,
            Limit: limit,
            Total: total,
        },
    }
    SendJSON(w, http.StatusOK, response)
}

func SendError(w http.ResponseWriter, err error) {
    statusCode := http.StatusInternalServerError
    errorDetail := &ErrorDetail{
        Code:    "INTERNAL_ERROR",
        Message: err.Error(),
    }

    if customErr, ok := err.(*errors.AppError); ok {
        statusCode = customErr.StatusCode
        errorDetail = &ErrorDetail{
            Code:    customErr.Code,
            Message: customErr.Message,
            Details: customErr.Details,
        }
    }

    response := Response{
        Success: false,
        Error:   errorDetail,
    }
    SendJSON(w, statusCode, response)
}

func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}
```

### Product Handlers (internal/controllers/barang_controller.go)
```go
package controllers

import (
    "ahadpos-go/internal/models"
    "ahadpos-go/internal/services"
    "ahadpos-go/internal/utils"
    "ahadpos-go/pkg/jwt"
    "net/http"
    "strconv"
    "github.com/go-chi/chi"
)

type BarangController struct {
    service *services.BarangService
}

func NewBarangController() *BarangController {
    return &BarangController{
        service: services.NewBarangService(),
    }
}

func (c *BarangController) RegisterRoutes(r chi.Router) {
    r.Get("/", c.List)
    r.Post("/", c.Create)
    r.Get("/{id}", c.GetByID)
    r.Put("/{id}", c.Update)
    r.Delete("/{id}", c.Delete)
    r.Get("/barcode/{barcode}", c.GetByBarcode)
}

func (c *BarangController) List(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    search := r.URL.Query().Get("search")
    statusStr := r.URL.Query().Get("status")

    // Set defaults
    if page == 0 {
        page = 1
    }
    if limit == 0 {
        limit = 20
    }

    var status *int
    if statusStr != "" {
        s, _ := strconv.Atoi(statusStr)
        status = &s
    }

    barang, total, err := c.service.List(page, limit, search, status)
    if err != nil {
        utils.SendError(w, err)
        return
    }

    utils.SendPaginated(w, barang, page, limit, total)
}

func (c *BarangController) GetByID(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
    if err != nil {
        utils.SendError(w, err)
        return
    }

    barang, err := c.service.GetByID(uint(id))
    if err != nil {
        utils.SendError(w, err)
        return
    }

    utils.SendSuccess(w, barang)
}

func (c *BarangController) Create(w http.ResponseWriter, r *http.Request) {
    var barang models.Barang
    if err := json.NewDecoder(r.Body).Decode(&barang); err != nil {
        utils.SendError(w, err)
        return
    }

    // Get user ID from JWT
    userID := jwt.GetUserIDFromContext(r.Context())
    barang.UpdatedBy = userID

    if err := c.service.Create(&barang); err != nil {
        utils.SendError(w, err)
        return
    }

    utils.SendCreated(w, barang)
}

func (c *BarangController) Update(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
    if err != nil {
        utils.SendError(w, err)
        return
    }

    var barang models.Barang
    if err := json.NewDecoder(r.Body).Decode(&barang); err != nil {
        utils.SendError(w, err)
        return
    }

    userID := jwt.GetUserIDFromContext(r.Context())
    barang.UpdatedBy = userID

    if err := c.service.Update(uint(id), &barang); err != nil {
        utils.SendError(w, err)
        return
    }

    utils.SendSuccess(w, barang)
}

func (c *BarangController) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
    if err != nil {
        utils.SendError(w, err)
        return
    }

    if err := c.service.Delete(uint(id)); err != nil {
        utils.SendError(w, err)
        return
    }

    utils.SendSuccess(w, map[string]string{"message": "barang deleted"})
}

func (c *BarangController) GetByBarcode(w http.ResponseWriter, r *http.Request) {
    barcode := chi.URLParam(r, "barcode")

    barang, err := c.service.GetByBarcode(barcode)
    if err != nil {
        utils.SendError(w, err)
        return
    }

    utils.SendSuccess(w, barang)
}
```

## 7. Middleware

### JWT Authentication Middleware (internal/middleware/auth.go)
```go
package middleware

import (
    "ahadpos-go/pkg/jwt"
    "net/http"
    "strings"
)

func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
        claims, err := jwt.ValidateToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Add user info to context
        ctx := jwt.SetUserIDContext(r.Context(), claims.UserID)
        ctx = jwt.SetUserRoleContext(ctx, claims.Role)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func RequireRole(roles ...string) func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userRole := jwt.GetUserRoleFromContext(r.Context())

            for _, role := range roles {
                if userRole == role {
                    next.ServeHTTP(w, r)
                    return
                }
            }

            http.Error(w, "Forbidden", http.StatusForbidden)
        })
    }
}
```

### CORS Middleware (internal/middleware/cors.go)
```go
package middleware

import (
    "github.com/go-chi/cors"
)

func CORS() func(http.Handler) http.Handler {
    return cors.Handler(cors.Options{
        AllowedOrigins:   []string{"*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300,
    })
}
```

### Logger Middleware (internal/middleware/logger.go)
```go
package middleware

import (
    "net/http"
    "time"
    "github.com/rs/zerolog/log"
)

func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Wrap response writer to capture status code
        ww := &responseWriter{ResponseWriter: w}

        next.ServeHTTP(ww, r)

        log.Info().
            Str("method", r.Method).
            Str("path", r.URL.Path).
            Str("query", r.URL.RawQuery).
            Int("status", ww.statusCode).
            Dur("duration", time.Since(start)).
            Msg("HTTP Request")
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

## 8. Main Application Setup

### Configuration (internal/config/config.go)
```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig      `mapstructure:"jwt"`
}

type ServerConfig struct {
    Port         int    `mapstructure:"port"`
    Mode         string `mapstructure:"mode"` // debug, release
    ReadTimeout  int    `mapstructure:"read_timeout"`
    WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Path string `mapstructure:"path"`
}

type JWTConfig struct {
    Secret     string `mapstructure:"secret"`
    ExpireTime int    `mapstructure:"expire_time"` // minutes
}

var AppConfig *Config

func Load(configPath string) error {
    viper.SetConfigFile(configPath)
    viper.SetConfigType("yaml")

    // Set defaults
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("server.mode", "debug")
    viper.SetDefault("server.read_timeout", 30)
    viper.SetDefault("server.write_timeout", 30)
    viper.SetDefault("database.path", "ahadpos.db")
    viper.SetDefault("jwt.expire_time", 15)

    if err := viper.ReadInConfig(); err != nil {
        return err
    }

    AppConfig = &Config{}
    return viper.Unmarshal(AppConfig)
}
```

### Main Entry Point (cmd/api/main.go)
```go
package main

import (
    "ahadpos-go/internal/config"
    "ahadpos-go/internal/controllers"
    "ahadpos-go/internal/middleware"
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
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

    // Connect to database
    if err := database.Connect(config.AppConfig.Database.Path); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer database.Close()

    // Create router
    r := chi.NewRouter()

    // Global middleware
    r.Use(middleware.RequestID)
    r.Use(middleware.RealIP)
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.Timeout(60 * time.Second))
    r.Use(middleware.CORS())
    r.Use(middleware.Logger)

    // Health check
    r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("OK"))
    })

    // API routes
    r.Route("/api/v1", func(r chi.Router) {
        // Public routes
        r.Route("/auth", func(r chi.Router) {
            authController := controllers.NewAuthController()
            authController.RegisterRoutes(r)
        })

        // Protected routes
        r.Group(func(r chi.Router) {
            r.Use(middleware.Auth)

            // User management
            r.Route("/users", func(r chi.Router) {
                userController := controllers.NewUserController()
                userController.RegisterRoutes(r)
            })

            // Products
            r.Route("/barang", func(r chi.Router) {
                barangController := controllers.NewBarangController()
                barangController.RegisterRoutes(r)
            })

            // Sales
            r.Route("/penjualan", func(r chi.Router) {
                penjualanController := controllers.NewPenjualanController()
                penjualanController.RegisterRoutes(r)
            })

            // POS
            r.Route("/pos", func(r chi.Router) {
                posController := controllers.NewPosController()
                posController.RegisterRoutes(r)
            })

            // Reports
            r.Route("/laporan", func(r chi.Router) {
                laporanController := controllers.NewLaporanController()
                laporanController.RegisterRoutes(r)
            })
        })
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

    // Graceful shutdown
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal().Err(err).Msg("Server failed to start")
        }
    }()

    // Wait for interrupt signal
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

## 9. Custom Errors

### Error Types (pkg/errors/errors.go)
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

func NewConflictError(message string) *AppError {
    return &AppError{
        Code:       "CONFLICT",
        Message:    message,
        StatusCode: http.StatusConflict,
    }
}

func NewAuthenticationError(message string) *AppError {
    return &AppError{
        Code:       "AUTHENTICATION_ERROR",
        Message:    message,
        StatusCode: http.StatusUnauthorized,
    }
}

func NewAuthorizationError(message string) *AppError {
    return &AppError{
        Code:       "AUTHORIZATION_ERROR",
        Message:    message,
        StatusCode: http.StatusForbidden,
    }
}
```

## 10. Docker Configuration

### Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
```

### docker-compose.yml
```yaml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    volumes:
      - ./data:/root/data
    environment:
      - SERVER_MODE=release
      - DB_PATH=/root/data/ahadpos.db
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
```

## 11. Testing Example

### Service Test (internal/services/barang_service_test.go)
```go
package services_test

import (
    "ahadpos-go/internal/models"
    "ahadpos-go/internal/repository"
    "ahadpos-go/internal/services"
    "ahadpos-go/pkg/database"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type BarangServiceTestSuite struct {
    suite.Suite
    db     *gorm.DB
    service *services.BarangService
}

func (suite *BarangServiceTestSuite) SetupTest() {
    var err error
    suite.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    assert.NoError(suite.T(), err)

    // Migrate tables
    suite.db.AutoMigrate(
        &models.Barang{},
        &models.KategoriBarang{},
        &models.SatuanBarang{},
        &models.RakBarang{},
    )

    // Set database connection
    database.DB = suite.db

    // Create service
    suite.service = services.NewBarangService()
}

func (suite *BarangServiceTestSuite) TearDownTest() {
    db, _ := suite.db.DB()
    db.Close()
}

func (suite *BarangServiceTestSuite) TestCreateBarang() {
    barang := &models.Barang{
        Barcode:  "123456789",
        Nama:     "Test Product",
        SatuanID: 1,
        Status:   1,
    }

    err := suite.service.Create(barang)
    assert.NoError(suite.T(), err)
    assert.NotZero(suite.T(), barang.ID)
}

func (suite *BarangServiceTestSuite) TestGetByBarcode() {
    barang := &models.Barang{
        Barcode:  "123456789",
        Nama:     "Test Product",
        SatuanID: 1,
        Status:   1,
    }
    suite.service.Create(barang)

    found, err := suite.service.GetByBarcode("123456789")
    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), barang.Nama, found.Nama)
}

func TestBarangServiceTestSuite(t *testing.T) {
    suite.Run(t, new(BarangServiceTestSuite))
}
```

## 12. API Documentation (Swagger)

### Swagger Annotations Example
```go
// @title AhadPOS 3 API
// @version 1.0
// @description REST API for AhadPOS 3 Point of Sale System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// ListBarangs godoc
// @Summary List all products
// @Description Get paginated list of products with optional filtering
// @Tags barang
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param search query string false "Search by barcode or name"
// @Param status query int false "Filter by status (0=inactive,1=active)"
// @Security BearerAuth
// @Success 200 {object} utils.Response{data=[]models.Barang,meta=utils.Meta}
// @Failure 401 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /barang [get]
```

This technical guide provides the foundation for implementing the Go-based REST API. Each section can be expanded based on specific requirements during the implementation phases.
