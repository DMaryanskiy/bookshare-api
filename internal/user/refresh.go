package user

import (
	"net/http"
	"time"

	"github.com/DMaryanskiy/bookshare-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Validates refresh token and returns a new access and refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token  body  RefreshRequest  true  "Refresh token to validate"
// @Success      200  {object}  map[string]string  "New access and refresh tokens"
// @Failure      400  {object}  map[string]string  "Missing refresh token"
// @Failure      401  {object}  map[string]string  "Invalid refresh token"
// @Failure      500  {object}  map[string]string  "Server error creating or deleting tokens"
// @Router       /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh token"})
		return
	}

	userID, err := h.TokenStore.VerifyRefreshToken(c, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	err = h.TokenStore.DeleteRefreshToken(c, req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete old refresh token"})
		return
	}
	newRefresh, err := h.TokenStore.CreateRefreshToken(c, userID, 24*time.Hour)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create refresh token"})
        return
    }

	accessToken, err := utils.GenerateAccessToken(userID, 15*time.Minute)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create access token"})
        return
    }

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefresh,
	})
}
