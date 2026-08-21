package handler

import (
	"errors"
	"net/http"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication and current-user HTTP requests.
type AuthHandler struct {
	userService *service.UserService
}

// NewAuthHandler creates an authentication handler backed by userService.
func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// RegisterRequest represents the payload for user registration.
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the payload for user authentication.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse represents public user data returned by the API.
type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// LoginResponse contains the access token issued after authentication.
type LoginResponse struct {
	Token string `json:"token"`
}

// Register creates a user account and returns its public data.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	user, err := h.userService.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}
	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	c.JSON(http.StatusCreated, response)
}

// Login authenticates a user and returns an access token.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	tokenString, err := h.userService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) ||
			errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Me returns the authenticated user's public data.
func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetInt64("userID")
	if userID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid userID"})
		return
	}
	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}
	response := UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
	c.JSON(http.StatusOK, response)
}
