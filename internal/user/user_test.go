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
	"github.com/DMaryanskiy/bookshare-api/internal/task/distributor"
	"github.com/DMaryanskiy/bookshare-api/internal/tests"
	"github.com/DMaryanskiy/bookshare-api/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type FakeTokenStore struct {
	Tokens map[string]string
}

func NewFakeTokenStore() *FakeTokenStore {
	return &FakeTokenStore{Tokens: make(map[string]string)}
}

func (f *FakeTokenStore) CreateRefreshToken(_ context.Context, userID string, ttl time.Duration) (string, error) {
	token := "fake-" + userID
	f.Tokens[token] = userID
	return token, nil
}

func (f *FakeTokenStore) VerifyRefreshToken(_ context.Context, token string) (string, error) {
	userID, ok := f.Tokens[token]
	if !ok {
		return "", redis.Nil
	}
	return userID, nil
}

func (f *FakeTokenStore) DeleteRefreshToken(_ context.Context, token string) error {
	delete(f.Tokens, token)
	return nil
}

func setupRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	tests.SetupTestDB(t)
	tests.SetupTestRedis()

	router := gin.Default()

	dist := distributor.NewTaskDistributor(tests.TestRedis.Options().Addr)
	tokenStore := &FakeTokenStore{} // can use mock for now

	h := user.NewHandler(dist, tokenStore)
	router.POST("/register", h.RegisterUser)
	router.GET("/verify", h.VerifyEmail)

	return router
}

func TestRegisterUser_Success(t *testing.T) {
	r := setupRouter(t)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "strongpass123",
	}
	data, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	var user models.User
	err := db.DB.Table("auth.users").Where("email = ?", "test@example.com").First(&user).Error
	require.NoError(t, err)
	require.False(t, user.IsVerified)
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	r := setupRouter(t)

	email := "dupe@example.com"
	_ = db.DB.Table("auth.users").Create(&models.User{
		Email:        email,
		PasswordHash: "somehash",
	}).Error

	body := map[string]string{
		"email":    email,
		"password": "12345678",
	}
	data, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusConflict, w.Code)
}

func TestRegisterUser_InvalidInput(t *testing.T) {
	r := setupRouter(t)

	body := map[string]string{
		"password": "12345678",
	}
	data, _ := json.Marshal(body)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerifyEmail_Success(t *testing.T) {
	r := setupRouter(t)

	user := models.User{
		Email:        "verify@example.com",
		PasswordHash: "somehash",
	}
	require.NoError(t, db.DB.Table("auth.users").Create(&user).Error)

	token := "valid-token"
	exp := time.Now().Add(10 * time.Minute)

	db.DB.Table("auth.verification_tokens").Create(&models.VerificationToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: exp,
	})

	req, _ := http.NewRequest("GET", "/verify?token="+token+"&uid="+user.ID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var updated models.User
	_ = db.DB.Table("auth.users").First(&updated, "id = ?", user.ID)
	require.True(t, updated.IsVerified)
}

func TestVerifyEmail_InvalidToken(t *testing.T) {
	r := setupRouter(t)

	user := models.User{
		Email:        "badtoken@example.com",
		PasswordHash: "somehash",
	}
	require.NoError(t, db.DB.Table("auth.users").Create(&user).Error)

	req, _ := http.NewRequest("GET", "/verify?token=invalid-token&uid="+user.ID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerifyEmail_ExpiredToken(t *testing.T) {
	r := setupRouter(t)

	user := models.User{
		Email:        "expired@example.com",
		PasswordHash: "somehash",
	}
	require.NoError(t, db.DB.Table("auth.users").Create(&user).Error)

	token := "expired-token"
	exp := time.Now().Add(-10 * time.Minute) // expired

	db.DB.Table("auth.verification_tokens").Create(&models.VerificationToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: exp,
	})

	req, _ := http.NewRequest("GET", "/verify?token="+token+"&uid="+user.ID.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
}
