package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid"`
	Action    string    `gorm:"not null"`
	Metadata  string    // optional JSON string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
