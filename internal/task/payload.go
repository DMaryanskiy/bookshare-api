package task

const TaskSendVerificationEmail = "send_verification_email"

type PayloadSendVerificationEmail struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
}
