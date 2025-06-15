package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/DMaryanskiy/bookshare-api/internal/middleware"
	"github.com/DMaryanskiy/bookshare-api/internal/task/distributor"
	"github.com/DMaryanskiy/bookshare-api/internal/tests"
	"github.com/DMaryanskiy/bookshare-api/internal/user"
	"github.com/DMaryanskiy/bookshare-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupAuthRouter(t *testing.T, tokenStore *FakeTokenStore) *gin.Engine {
	gin.SetMode(gin.TestMode)

	tests.SetupTestDB(t)
	tests.SetupTestRedis()

	router := gin.Default()
	dist := distributor.NewTaskDistributor(tests.TestRedis.Options().Addr)

	h := user.NewHandler(dist, tokenStore)
	router.POST("/login", h.LoginUser)
	router.POST("/refresh", h.RefreshToken)

	protected := router.Group("/protected")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	return router
}

func TestLogin_Success(t *testing.T) {
	tokenStore := NewFakeTokenStore()
	r := setupAuthRouter(t, tokenStore)

	// Create a verified user first
	hashedPassword, _ := utils.HashPassword("password123")
	user := models.User{
		Email:        "login@example.com",
		PasswordHash: hashedPassword,
		IsVerified:   true,
	}
	require.NoError(t, db.DB.Table("auth.users").Create(&user).Error)

	body := map[string]string{
		"email":    "login@example.com",
		"password": "password123",
	}
	data, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Contains(t, resp, "access_token")
	require.Contains(t, resp, "refresh_token")
}

func TestLogin_Failure(t *testing.T) {
	tokenStore := NewFakeTokenStore()
	r := setupAuthRouter(t, tokenStore)

	body := map[string]string{
		"email":    "nonexistent@example.com",
		"password": "badpassword",
	}
	data, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRefreshToken_Success(t *testing.T) {
	tokenStore := NewFakeTokenStore() // shared instance
	r := setupAuthRouter(t, tokenStore)

	// Use FakeTokenStore to create a token

	hashedPassword, _ := utils.HashPassword("password123")
	user := models.User{
		Email:        "refresh@example.com",
		PasswordHash: hashedPassword,
		IsVerified:   true,
	}
	require.NoError(t, db.DB.Table("auth.users").Create(&user).Error)
	refreshToken, err := tokenStore.CreateRefreshToken(context.TODO(), user.ID.String(), time.Hour*24)
	require.NoError(t, err)

	body := map[string]string{
		"refresh_token": refreshToken,
	}
	data, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Contains(t, resp, "access_token")
	require.Contains(t, resp, "refresh_token")
}

func TestProtectedRoute_Access(t *testing.T) {
	tokenStore := NewFakeTokenStore()
	r := setupAuthRouter(t, tokenStore)

	// Create a verified user and login to get token
	hashedPassword, _ := utils.HashPassword("password123")
	user := models.User{
		Email:        "protected@example.com",
		PasswordHash: hashedPassword,
		IsVerified:   true,
	}
	require.NoError(t, db.DB.Table("auth.users").Create(&user).Error)

	// Manually generate JWT token (or login and get it)
	token, err := utils.GenerateAccessToken(user.ID.String(), time.Hour)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestProtectedRoute_AccessDenied(t *testing.T) {
	tokenStore := NewFakeTokenStore()
	r := setupAuthRouter(t, tokenStore)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}
