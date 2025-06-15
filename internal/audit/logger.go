package audit

import (
	"context"
	"encoding/json"
	"log"

	"github.com/DMaryanskiy/bookshare-api/internal/db"
	"github.com/DMaryanskiy/bookshare-api/internal/db/models"
	"github.com/google/uuid"
)

func Log(userID uuid.UUID, action string, meta any) {
	metaStr := ""
	if meta != nil {
		if b, err := json.Marshal(meta); err == nil {
			metaStr = string(b)
		}
	}

	entry := models.AuditLog{
		UserID:   userID,
		Action:   action,
		Metadata: metaStr,
	}

	if err := db.DB.WithContext(context.Background()).Table("logs.audit_logs").Create(&entry).Error; err != nil {
		log.Printf("Failed to write audit log: %v", err)
	}
}
