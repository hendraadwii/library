package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hendraadwii/library/internal/models"
)

type UserHandler struct {
	UserService *models.UserService
}

func NewUserHandler(userService *models.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// GetUsers returns a paginated list of users with optional filters
// @Summary Get users
// @Description Get a paginated list of users with optional filters
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10)"
// @Param search query string false "Search by name or email"
// @Param role query string false "Filter by role (admin/member)"
// @Param status query string false "Filter by status (active/locked)"
// @Success 200 {object} PaginatedUsersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	// Calculate offset
	offset := (page - 1) * perPage

	// Parse filter parameters
	search := c.Query("search")
	role := models.UserRole(c.Query("role"))
	status := models.UserStatus(c.Query("status"))

	// Get users from database with filters
	users, err := h.UserService.List(perPage, offset, search, role, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting users: " + err.Error()})
		return
	}

	// Get total count for pagination with filters
	totalCount, err := h.UserService.Count(search, role, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error counting users: " + err.Error()})
		return
	}

	// Calculate total pages
	totalPages := totalCount / perPage
	if totalCount%perPage > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"total_count": totalCount,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

// CreateUser creates a new user
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User details"
// @Success 201 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "Email already exists"
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Email    string          `json:"email" binding:"required,email"`
		Password string          `json:"password" binding:"required,min=6"`
		FullName string          `json:"full_name" binding:"required,min=2"`
		Role     models.UserRole `json:"role" binding:"required,oneof=admin member"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Validate email format
	if !isValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validate password strength
	if !isValidPassword(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long and contain at least one number"})
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
		Status:   models.StatusActive, // Explicitly set status
	}

	if err := h.UserService.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user: " + err.Error()})
		return
	}

	// Clear sensitive data before sending response
	user.Password = ""
	user.PasswordHash = ""

	c.JSON(http.StatusCreated, user)
}

// UpdateUser updates an existing user
// @Summary Update user
// @Description Update an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "User details"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Parse user ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Email    string          `json:"email" binding:"omitempty,email"`
		FullName string          `json:"full_name" binding:"omitempty,min=2"`
		Role     models.UserRole `json:"role" binding:"omitempty,oneof=admin member"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing user
	user, err := h.UserService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// If email is being updated, check if it's already taken
	if req.Email != "" && req.Email != user.Email {
		if !isValidEmail(req.Email) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
			return
		}

		existingUser, err := h.UserService.GetByEmail(req.Email)
		if err == nil && existingUser.ID != user.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}
		user.Email = req.Email
	}

	// Update other fields if provided
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := h.UserService.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user: " + err.Error()})
		return
	}

	// Clear sensitive data before sending response
	user.Password = ""
	user.PasswordHash = ""

	c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Parse user ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.UserService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ResetUserPassword resets a user's password
// @Summary Reset user password
// @Description Reset a user's password
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param password body ResetPasswordRequest true "New password"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/{id}/reset-password [post]
func (h *UserHandler) ResetUserPassword(c *gin.Context) {
	// Parse user ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing user
	user, err := h.UserService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Update password
	user.Password = req.Password
	if err := h.UserService.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error resetting password: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

// ToggleUserStatus toggles a user's status (lock/unlock)
// @Summary Toggle user status
// @Description Toggle a user's status between active and locked
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param status body ToggleStatusRequest true "New status"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/{id}/status [put]
func (h *UserHandler) ToggleUserStatus(c *gin.Context) {
	// Parse user ID
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=active locked"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get existing user
	user, err := h.UserService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Update status
	user.Status = models.UserStatus(req.Status)
	if err := h.UserService.Update(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserCount returns the total number of users (admin only)
// @Summary Get user count
// @Description Get the total number of users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} CountResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /users/count [get]
func (h *UserHandler) GetUserCount(c *gin.Context) {
	count, err := h.UserService.Count("", "", "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user count"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_users": count})
}

// Helper functions for validation
func isValidEmail(email string) bool {
	// Basic email validation - you might want to use a more robust validation
	return len(email) >= 3 && len(email) <= 255 && strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isValidPassword(password string) bool {
	// Password must be at least 6 characters and contain at least one number
	if len(password) < 6 {
		return false
	}
	hasNumber := false
	for _, c := range password {
		if c >= '0' && c <= '9' {
			hasNumber = true
			break
		}
	}
	return hasNumber
}
