package user

import (
	"net/http"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/gin-gonic/gin"
)

// GetMe godoc
// @Summary      Get current user
// @Description  Returns the authenticated user's details
// @Tags         user
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "User info"
// @Failure      404  {object}  map[string]string  "User not found"
// @Failure      500  {object}  map[string]string  "User ID missing from context"
// @Router       /users/me [get]
func (h *Handler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user ID missing from context"})
		return
	}

	var user models.User
	if err := db.DB.Table("auth.users").First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"verified":   user.IsVerified,
		"created_at": user.CreatedAt,
	})
}
