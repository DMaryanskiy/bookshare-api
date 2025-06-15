package user

import (
	"net/http"
	"time"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/DMaryanskiy/bookshare-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginUser godoc
// @Summary      Login user
// @Description  Authenticates a user and returns access and refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true  "Email and Password"
// @Success      200  {object}  map[string]string  "Returns access and refresh tokens"
// @Failure      400  {object}  map[string]string  "Invalid input"
// @Failure      401  {object}  map[string]string  "Unauthorized - invalid credentials or not verified"
// @Failure      500  {object}  map[string]string  "Internal error"
func (h *Handler) LoginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login input"})
		return
	}

	var user models.User
	if err := db.DB.Table("auth.users").
		Where("email = ?", req.Email).
		First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	ok, err := utils.CheckPasswordHash(user.PasswordHash, req.Password)
	if !ok || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
        return
	}

	if !user.IsVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not verified"})
	}

	accessToken, err := utils.GenerateAccessToken(user.ID.String(), 15 * time.Minute)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create access token"})
        return
    }

	refreshToken, err := h.TokenStore.CreateRefreshToken(c, user.ID.String(), 24 * time.Hour)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create refresh token"})
        return
    }

	c.JSON(http.StatusOK, gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}
