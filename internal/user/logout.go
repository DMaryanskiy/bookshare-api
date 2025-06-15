package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// LogoutUser godoc
// @Summary      Logout user
// @Description  Revokes the refresh token and logs the user out
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token  body  LogoutRequest  true  "Refresh token to revoke"
// @Success      200  {object}  map[string]string  "Successfully logged out"
// @Failure      400  {object}  map[string]string  "Refresh token required"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
		return
	}

	err := h.TokenStore.DeleteRefreshToken(c, req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not revoke refresh token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}