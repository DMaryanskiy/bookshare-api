package main

import (
	"log"
	"os"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/email"
	"github.com/DMaryanskiy/bookshare-api/internal/task/processor"
	"github.com/joho/godotenv"
)

func main() {
	debug_mode := os.Getenv("DEBUG")
	if debug_mode != "" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Failed to load env:", err)
		}
	}

	db.InitDB()

	redisAddr := os.Getenv("REDIS_ADDR")
	sender := email.NewEmailSender()

	taskProcessor := processor.NewTaskProcessor(sender)
	if err := taskProcessor.Start(redisAddr); err != nil {
		log.Fatal("failed to start worker:", err)
	}
}
