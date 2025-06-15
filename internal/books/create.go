package books

import (
	"net/http"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateBook godoc
// @Summary      Create a new book
// @Description  Creates a new book record for the authenticated user
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body  BookInput  true  "Book details"
// @Success      201  {object}  models.Book  "Book created successfully"
// @Failure      400  {object}  map[string]string  "Invalid input"
// @Failure      500  {object}  map[string]string  "Failed to create book"
// @Router       /books [post]
func (h *Handler) CreateBook(c *gin.Context) {
	userID := c.GetString("user_id")

	var req BookInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	book := models.Book{
		UserID:      uuid.MustParse(userID),
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
	}

	if err := db.DB.Table("books.books").Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create book"})
		return
	}

	c.JSON(http.StatusCreated, book)
}
