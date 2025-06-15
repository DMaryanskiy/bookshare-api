package user

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DMaryanskiy/bookshare-api/internal/audit"
	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/DMaryanskiy/bookshare-api/internal/task"
	"github.com/DMaryanskiy/bookshare-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Creates a new user and sends a verification email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body  RegisterUserRequest  true  "User credentials"
// @Success      201  {object}  map[string]string  "Registration successful"
// @Failure      400  {object}  map[string]string  "Invalid request"
// @Failure      409  {object}  map[string]string  "Email already exists"
// @Failure      500  {object}  map[string]string  "Server error"
// @Router       /auth/register [post]
func (h *Handler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		errStr := fmt.Sprintf("could not hash password: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errStr})
		return
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}
	if err := db.DB.Table("auth.users").Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	payload := task.PayloadSendVerificationEmail{
		UserId: user.ID.String(),
		Email:  user.Email,
	}
	if err := h.TaskDistributor.DistributeVerificationEmail(context.Background(), payload); err != nil {
		errStr := fmt.Sprintf("failed to enqueue email: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errStr})
		return
	}

	audit.Log(user.ID, "registration_success", map[string]string{
		"ip":         c.ClientIP(),
		"user_agent": c.GetHeader("User-Agent"),
	})

	c.JSON(http.StatusCreated, gin.H{"message": "registration successful, verification email sent"})
}
