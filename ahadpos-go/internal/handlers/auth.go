package handlers

import (
    "net/http"

    "github.com/ahadpos/go/internal/config"
    "github.com/ahadpos/go/internal/models"
    "github.com/ahadpos/go/pkg/database"
    "github.com/ahadpos/go/pkg/utils"
    "github.com/gin-gonic/gin"
)

var db = database.DB

type AuthHandler struct {
    cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
    return &AuthHandler{cfg: cfg}
}

// LoginRequest represents the login request body
type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

// LoginResponse represents the login response
type LoginResponse struct {
    Token    string `json:"token"`
    UserID   uint   `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    Name     string `json:"name"`
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user models.User
    if err := db.Where("username = ? AND status = 1", req.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    if !utils.CheckPassword(user.Password, req.Password) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := utils.GenerateJWT(user.ID, user.Username, user.Role, h.cfg.JWT.Expiration)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    c.JSON(http.StatusOK, LoginResponse{
        Token:    token,
        UserID:   user.ID,
        Username: user.Username,
        Role:     user.Role,
        Name:     user.Name,
    })
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
    userID, _ := c.Get("user_id")

    var user models.User
    if err := db.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}
