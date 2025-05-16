package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hendraadwii/library/internal/auth"
	"github.com/hendraadwii/library/internal/models"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	UserService  *models.UserService
	TokenManager *auth.TokenManager
	RateLimiter  *auth.RateLimiter
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email    string          `json:"email" binding:"required,email"`
	Password string          `json:"password" binding:"required,min=6"`
	FullName string          `json:"full_name" binding:"required"`
	Role     models.UserRole `json:"role" binding:"omitempty"`
}

// RefreshRequest represents a refresh token request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userService *models.UserService, tokenManager *auth.TokenManager, rateLimiter *auth.RateLimiter) *AuthHandler {
	return &AuthHandler{
		UserService:  userService,
		TokenManager: tokenManager,
		RateLimiter:  rateLimiter,
	}
}

// Login handles user login
// @Summary Login user
// @Description Authenticate a user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} auth.TokenPair
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 429 {object} ErrorResponse "Too many requests"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Authenticate user
	user, err := h.UserService.Authenticate(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate tokens
	tokens, err := h.TokenManager.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating tokens"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Register handles user registration
// @Summary Register user
// @Description Register a new user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} auth.TokenPair
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "Email already exists"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Check if email already exists
	_, err := h.UserService.GetByEmail(req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	}

	// Create new user
	user := &models.User{
		Email:    req.Email,
		Password: req.Password, // Will be hashed in the service
		FullName: req.FullName,
		Role:     req.Role,
	}

	// Default role if not specified
	if user.Role == "" {
		user.Role = models.RoleMember
	}

	// Admin registration might require additional checks
	if user.Role == models.RoleAdmin {
		// In a real application, you might want to check if the request is authorized to create an admin
		// For now, let's prevent admin registration through the API
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin creation not allowed through registration"})
		return
	}

	// Save user to database
	if err := h.UserService.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user: " + err.Error()})
		return
	}

	// Generate tokens
	tokens, err := h.TokenManager.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating tokens"})
		return
	}

	c.JSON(http.StatusCreated, tokens)
}

// RefreshToken handles token refresh
// @Summary Refresh tokens
// @Description Refresh the access token using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token request"
// @Success 200 {object} auth.TokenPair
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse "Invalid or expired refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Refresh the tokens
	tokens, err := h.TokenManager.RefreshTokens(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	c.JSON(http.StatusOK, tokens)
}
