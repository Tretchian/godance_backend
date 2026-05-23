package models

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	ParticipantID uint      `gorm:"not null;index"`
	Participant   User      `gorm:"foreignKey:ParticipantID"`
	StorageKey    *string   `gorm:"type:varchar"` // NULL если deleted
	UploadedAt    time.Time `gorm:"autoCreateTime"`
	ExpiresAt     time.Time `gorm:"type:timestamptz"`
}
