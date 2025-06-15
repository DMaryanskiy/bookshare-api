package email

import (
	"log"
	"os"
	"strconv"

	"gopkg.in/mail.v2"
)

type EmailSender struct {
	From string
	Host string
	Port int
	User string
	Pass string
}

func NewEmailSender() *EmailSender {
	portStr := os.Getenv("SMTP_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalln("failed to convert port to int:", err)
	}

	return &EmailSender{
		From: os.Getenv("SMTP_USER"),
		Host: os.Getenv("SMTP_HOST"),
		Port: port,
		User: os.Getenv("SMTP_USER"),
		Pass: os.Getenv("SMTP_PASS"),
	}
}

func (s *EmailSender) Send(to, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := mail.NewDialer(s.Host, s.Port, s.User, s.Pass)

	return d.DialAndSend(m)
}
