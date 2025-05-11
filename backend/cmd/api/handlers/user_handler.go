package handlers

import (
	"net/http"

	"github.com/aksaaaraa/library/internal/models"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *models.UserService
}

func NewUserHandler(userService *models.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

// GetUserCount returns the total number of users (admin only)
func (h *UserHandler) GetUserCount(c *gin.Context) {
	count, err := h.UserService.Count()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user count"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_users": count})
}
