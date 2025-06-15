package books

import (
	"net/http"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/gin-gonic/gin"
)

// DeleteBook godoc
// @Summary      Delete a book
// @Description  Deletes a book owned by the authenticated user
// @Tags         books
// @Produce      json
// @Param        id   path      string  true  "Book ID"
// @Success      200  {object}  map[string]string  "Book deleted successfully"
// @Failure      500  {object}  map[string]string  "Failed to delete book"
// @Router       /books/{id} [delete]
func (h *Handler) DeleteBook(c *gin.Context) {
	userID := c.GetString("user_id")
	bookID := c.Param("id")

	if err := db.DB.Table("books.books").
		Where("id = ? AND user_id = ?", bookID, userID).
		Delete(&models.Book{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "book deleted"})
}
