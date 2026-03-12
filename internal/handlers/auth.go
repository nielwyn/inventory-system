package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nielwyn/inventory-system/internal/models"
	"github.com/nielwyn/inventory-system/internal/service"
	"github.com/nielwyn/inventory-system/pkg/logger"
	"github.com/nielwyn/inventory-system/pkg/response"
	"github.com/nielwyn/inventory-system/pkg/validator"
	"go.uber.org/zap"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, validator.FormatValidationError(err))
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		logger.Error("Registration failed", zap.Error(err))
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", gin.H{
		"user": user,
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, validator.FormatValidationError(err))
		return
	}

	loginResponse, err := h.authService.Login(&req)
	if err != nil {
		logger.Error("Login failed", zap.Error(err))
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Login successful", loginResponse)
}
