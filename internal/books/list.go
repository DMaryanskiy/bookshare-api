package books

import (
	"net/http"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/gin-gonic/gin"
)

// ListBooks godoc
// @Summary      List all books
// @Description  Returns a list of all books owned by the authenticated user
// @Tags         books
// @Produce      json
// @Success      200  {array}   models.Book  "List of books"
// @Failure      500  {object}  map[string]string  "Could not fetch books"
// @Router       /books [get]
func (h *Handler) ListBooks(c *gin.Context) {
	userID := c.GetString("user_id")
	var books []models.Book

	if err := db.DB.Table("books.books").
		Where("user_id = ?", userID).
		Find(&books).
		Order("created_at desc").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch books"})
		return
	}

	c.JSON(http.StatusOK, books)
}
