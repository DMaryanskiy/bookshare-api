package books

import (
	"net/http"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/gin-gonic/gin"
)

// GetBook godoc
// @Summary      Get a book
// @Description  Retrieves a specific book owned by the authenticated user
// @Tags         books
// @Produce      json
// @Param        id   path      string  true  "Book ID"
// @Success      200  {object}  models.Book  "Book retrieved successfully"
// @Failure      404  {object}  map[string]string  "Book not found"
// @Router       /books/{id} [get]
func (h *Handler) GetBook(c *gin.Context) {
	userID := c.GetString("user_id")
	bookID := c.Param("id")

	var book models.Book
	if err := db.DB.Table("books.books").
		Where("id = ? AND user_id = ?", bookID, userID).
		First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	c.JSON(http.StatusOK, book)
}
