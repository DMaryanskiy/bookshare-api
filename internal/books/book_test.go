package books_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DMaryanskiy/bookshare-api/internal/books"
	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/DMaryanskiy/bookshare-api/internal/middleware"
	"github.com/DMaryanskiy/bookshare-api/internal/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupBookRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	// Middleware for JWT auth
	router.Use(middleware.JWTAuthMiddleware())

	bh := books.NewHandler()
	router.POST("/books", bh.CreateBook)
	router.GET("/books/:id", bh.GetBook)
	router.PUT("/books/:id", bh.UpdateBook)
	router.DELETE("/books/:id", bh.DeleteBook)

	return router
}

func TestCreateBook_Success(t *testing.T) {
	tests.SetupTestDB(t)
	tests.SetupTestRedis()

	_, token := tests.CreateTestUser(t, "bookuser@example.com", "password123")
	r := setupBookRouter()

	body := map[string]string{
		"title":       "Test Book",
		"author":      "Author Name",
		"description": "Description",
	}
	data, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tests.GetAuthHeader(token))

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var resp models.Book
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, "Test Book", resp.Title)
	require.Equal(t, "Author Name", resp.Author)
	require.Equal(t, "Description", resp.Description)
}

func TestGetBook_Success(t *testing.T) {
	tests.SetupTestDB(t)
	tests.SetupTestRedis()

	user, token := tests.CreateTestUser(t, "reader@example.com", "password123")
	book := models.Book{
		Title:       "Readable Book",
		Author:      "Read Author",
		Description: "Description",
		UserID:      user.ID,
	}
	require.NoError(t, db.DB.Table("books.books").Create(&book).Error)

	r := setupBookRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/"+book.ID.String(), nil)
	req.Header.Set("Authorization", tests.GetAuthHeader(token))

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateBook_Success(t *testing.T) {
	tests.SetupTestDB(t)
	tests.SetupTestRedis()

	user, token := tests.CreateTestUser(t, "updater@example.com", "password123")
	book := models.Book{
		Title:       "Old Title",
		Author:      "Old Author",
		Description: "Description",
		UserID:      user.ID,
	}
	require.NoError(t, db.DB.Table("books.books").Create(&book).Error)

	r := setupBookRouter()

	updateBody := map[string]string{
		"title":  "Updated Title",
		"author": "Updated Author",
	}
	data, _ := json.Marshal(updateBody)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/books/"+book.ID.String(), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tests.GetAuthHeader(token))

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	// Fetch updated book
	var updated models.Book
	err := db.DB.Table("books.books").First(&updated, "id = ?", book.ID).Error
	require.NoError(t, err)
	require.Equal(t, "Updated Title", updated.Title)
}

func TestDeleteBook_Success(t *testing.T) {
	tests.SetupTestDB(t)
	tests.SetupTestRedis()

	user, token := tests.CreateTestUser(t, "deleter@example.com", "password123")
	book := models.Book{
		Title:       "Book to Delete",
		Author:      "Ghost Author",
		Description: "Description",
		UserID:      user.ID,
	}
	require.NoError(t, db.DB.Table("books.books").Create(&book).Error)

	r := setupBookRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/books/"+book.ID.String(), nil)
	req.Header.Set("Authorization", tests.GetAuthHeader(token))

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var count int64
	db.DB.Table("books.books").Model(&models.Book{}).Where("id = ?", book.ID).Count(&count)
	require.Equal(t, int64(0), count)
}

func TestCreateBook_Unauthorized(t *testing.T) {
	tests.SetupTestDB(t)
	tests.SetupTestRedis()

	r := setupBookRouter()

	body := map[string]string{
		"title":        "Unauthorized Book",
		"author":       "Nobody",
		"descriptions": "Description",
	}
	data, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	// No Authorization header

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusUnauthorized, w.Code)
}
