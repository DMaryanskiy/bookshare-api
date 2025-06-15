package books

import (
	"net/http"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/gin-gonic/gin"
)

// UpdateBook godoc
// @Summary      Update a book
// @Description  Updates the details of a book owned by the authenticated user
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      string     true  "Book ID"
// @Param        book  body      BookInput  true  "Updated book data"
// @Success      200  {object}  models.Book  "Updated book"
// @Failure      400  {object}  map[string]string  "Invalid input"
// @Failure      404  {object}  map[string]string  "Book not found"
// @Failure      500  {object}  map[string]string  "Failed to update book"
// @Router       /books/{id} [put]
func (h *Handler) UpdateBook(c *gin.Context) {
	userID := c.GetString("user_id")
	bookID := c.Param("id")

	var book models.Book
	if err := db.DB.Table("books.books").
		Where("id = ? AND user_id = ?", bookID, userID).
		First(&book).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}

	var req BookInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	book.Title = req.Title
	book.Author = req.Author
	book.Description = req.Description

	if err := db.DB.Table("books.books").Save(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update book"})
		return
	}

	c.JSON(http.StatusOK, book)
}
