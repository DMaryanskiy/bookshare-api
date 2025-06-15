package admin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DMaryanskiy/bookshare-api/internal/admin"
	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/middleware"
	"github.com/DMaryanskiy/bookshare-api/internal/tests"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupAdminRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	// Middleware for JWT auth
	router.Use(middleware.JWTAuthMiddleware(), middleware.AdminOnly())

	adm := admin.NewHandler()
	router.GET("/admin/users", adm.ListUsers)

	return router
}

func TestAdminListUsers_Success(t *testing.T) {
	tests.SetupTestDB(t)
	tests.SetupTestRedis()

	admin, token := tests.CreateTestUser(t, "admin@example.com", "adminpass")
	admin.Role = "admin"
	require.NoError(t, db.DB.Table("auth.users").Save(&admin).Error)

	r := setupAdminRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/users", nil)
	req.Header.Set("Authorization", tests.GetAuthHeader(token))

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestAdminAccess_ForbiddenForNonAdmin(t *testing.T) {
	tests.SetupTestDB(t)
	tests.SetupTestRedis()

	_, token := tests.CreateTestUser(t, "user2@example.com", "userpass")

	r := setupAdminRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/admin/users", nil)
	req.Header.Set("Authorization", tests.GetAuthHeader(token))

	r.ServeHTTP(w, req)
	require.Equal(t, http.StatusForbidden, w.Code)
}
