package models

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	ParticipantID uint      `gorm:"not null;index"`
	Participant   User      `gorm:"foreignKey:ParticipantID"`
	StorageKey    *string   `gorm:"type:varchar"`                                // NULL если deleted
	Status        string    `gorm:"type:varchar(20);not null;default:uploading"` // uploading | active | deleted
	ContentType   string    `gorm:"type:varchar(50)"`
	SizeBytes     int64     `gorm:"type:bigint"`
	UploadedAt    time.Time `gorm:"autoCreateTime"`
	ExpiresAt     time.Time `gorm:"type:timestamptz"`
}
