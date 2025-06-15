package admin

import (
	"net/http"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/gin-gonic/gin"
)

// ListUsers godoc
// @Summary      List all users
// @Description  Retrieves a list of all users ordered by creation date descending
// @Tags         admin
// @Produce      json
// @Success      200  {array}   models.User  "List of users"
// @Failure      500  {object}  map[string]string  "Could not retrieve users"
// @Router       /admin/users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	var users []models.User

	if err := db.DB.Table("auth.users").Order("created_at desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
