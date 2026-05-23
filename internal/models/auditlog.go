package models

import (
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	gorm.Model
	UserID     *uint `gorm:"index:idx_audit_logs_occurred_at"`
	User       *User
	Action     *string   `gorm:"type:varchar(100)"`
	Details    *[]byte   `gorm:"type:jsonb"`
	OccurredAt time.Time `gorm:"autoCreateTime;index:idx_audit_logs_occurred_at"`
}
