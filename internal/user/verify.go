package user

import (
	"net/http"
	"time"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/gin-gonic/gin"
)

type VerifyRequest struct {
	Token string `form:"token" binding:"required"`
	UID   string `form:"uid" binding:"required"`
}

// VerifyEmail godoc
// @Summary      Verify user email
// @Description  Verifies a user's email using the token and UID
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        token  query  string  true  "Verification token"
// @Param        uid    query  string  true  "User ID"
// @Success      200  {object}  map[string]string  "Email verified"
// @Failure      400  {object}  map[string]string  "Invalid or expired token"
// @Failure      500  {object}  map[string]string  "Verification failed due to server error"
// @Router       /auth/verify [get]
func (h *Handler) VerifyEmail(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token or user ID"})
		return
	}

	var token models.VerificationToken
	if err := db.DB.
		Table("auth.verification_tokens").
		Where("token = ? AND user_id = ?", req.Token, req.UID).
		First(&token).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification link"})
        return
	}

	if time.Now().After(token.ExpiresAt) {
		db.DB.Table("auth.verification_tokens").Delete(&token)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification link expired"})
        return
	}

	if err := db.DB.Table("auth.users").
		Model(&models.User{}).
		Where("id = ?", req.UID).
		Update("is_verified", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify user"})
		return
	}

	db.DB.Table("auth.verification_tokens").Delete(&token)

    c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
