package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/DMaryanskiy/bookshare-api/internal/email"
	"github.com/DMaryanskiy/bookshare-api/internal/task"
	"github.com/DMaryanskiy/bookshare-api/pkg/utils"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type TaskProcessor struct {
	EmailSender *email.EmailSender
}

func NewTaskProcessor(sender *email.EmailSender) *TaskProcessor {
	return &TaskProcessor{EmailSender: sender}
}

func (p *TaskProcessor) Start(redisAddr string) error {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 5},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(task.TaskSendVerificationEmail, p.handleSendVerificationEmail)

	log.Println("Email worker is running...")
	return srv.Run(mux)
}

func (p *TaskProcessor) handleSendVerificationEmail(ctx context.Context, t *asynq.Task) error {
	var payload task.PayloadSendVerificationEmail
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("invalid payload: %v", err)
	}

	token, err := utils.GenerateRandomToken(32)
	if err != nil {
		return fmt.Errorf("failed to generate token: %v", err)
	}
	expires := time.Now().Add(30 * time.Minute)
	db.DB.Table("auth.verification_tokens").Create(&models.VerificationToken{
		UserID: uuid.MustParse(payload.UserId),
		Token: token,
		ExpiresAt: expires,
	})

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	verificationLink := fmt.Sprintf("http://%s:8080/api/v1/verify?token=%s&uid=%s", host, token, payload.UserId)

	emailBody := fmt.Sprintf(`
        <h1>Verify your email</h1>
        <p>Click <a href="%s">here</a> to verify your account.</p>
        <p>If you did not request this, please ignore.</p>
    `, verificationLink)

	err = p.EmailSender.Send(payload.Email, "Verify your BookShare account", emailBody)
	if err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }

	log.Printf("Sent verification email to %s", payload.Email)
	return nil
}
