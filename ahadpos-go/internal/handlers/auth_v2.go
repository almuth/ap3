package handlers

import (
	"ahadpos-go/internal/services"
	"ahadpos-go/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req services.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.NewErrorResponse(utils.ErrValidationFailed))
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		c.JSON(utils.HTTPStatus(err), utils.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(response))
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.NewErrorResponse(utils.ErrUnauthorized))
		return
	}

	user, err := h.authService.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.NewErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, utils.NewSuccessResponse(user))
}
