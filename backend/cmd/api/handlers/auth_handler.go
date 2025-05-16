package handlers

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hendraadwii/library/internal/auth"
	"github.com/hendraadwii/library/internal/models"
)

// AuthHandler handles authentication-related requests
type AuthHandler struct {
	UserService          *models.UserService
	TokenManager         *auth.TokenManager
	RateLimiter          *auth.RateLimiter
	PasswordResetService *models.PasswordResetService
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

// ForgotPasswordRequest for forgot password endpoint
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyPinRequest for verify pin endpoint
type VerifyPinRequest struct {
	Email string `json:"email" binding:"required,email"`
	Pin   string `json:"pin" binding:"required"`
}

// ResetPasswordRequest for reset password endpoint
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Pin         string `json:"pin" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
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
		if err.Error() == "account is locked" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Your account has been locked. Please contact the administrator."})
			return
		}
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

// sendEmail sends an email with the given subject and body to the specified recipient
func sendEmail(to, subject, body string) error {
	// SMTP configuration from environment variables
	from := os.Getenv("SMTP_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		return fmt.Errorf("SMTP configuration missing. Please check environment variables: SMTP_FROM, SMTP_PASSWORD, SMTP_HOST, SMTP_PORT")
	}

	// Compose email
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s", from, to, subject, body)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %v. Please check your SMTP configuration and make sure you're using an App Password for Gmail", err)
	}

	// Log success (for debugging)
	fmt.Printf("Email sent successfully to: %s\n", to)
	return nil
}

// ForgotPassword handler
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Generate 6-digit PIN
	pin := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Set expired time (1 menit dari sekarang)
	expiredAt := time.Now().Add(1 * time.Minute)

	// Save PIN to database
	pinData := &models.PasswordResetPin{
		Email:     req.Email,
		Pin:       pin,
		ExpiredAt: expiredAt,
	}
	if err := h.PasswordResetService.Create(pinData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save PIN: " + err.Error()})
		return
	}

	// Return PIN in response (for development/testing)
	c.JSON(http.StatusOK, gin.H{
		"message":    "PIN berhasil dibuat",
		"pin":        pin,
		"expired_at": expiredAt.Format(time.RFC3339),
	})
}

// VerifyPin handler
func (h *AuthHandler) VerifyPin(c *gin.Context) {
	var req VerifyPinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	pinData, err := h.PasswordResetService.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PIN tidak ditemukan"})
		return
	}
	if pinData.Pin != req.Pin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN salah"})
		return
	}
	if time.Now().After(pinData.ExpiredAt) {
		h.PasswordResetService.DeleteByEmail(req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN sudah expired"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "PIN valid"})
}

// ResetPassword handler
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}
	pinData, err := h.PasswordResetService.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PIN tidak ditemukan"})
		return
	}
	if pinData.Pin != req.Pin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN salah"})
		return
	}
	if time.Now().After(pinData.ExpiredAt) {
		h.PasswordResetService.DeleteByEmail(req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "PIN sudah expired"})
		return
	}
	// Update password user
	user, err := h.UserService.GetByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	user.Password = req.NewPassword
	err = h.UserService.Update(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal update password: " + err.Error()})
		return
	}
	// Hapus PIN
	h.PasswordResetService.DeleteByEmail(req.Email)
	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil direset."})
}
