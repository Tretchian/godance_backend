package models

import (
	"time"

	"gorm.io/gorm"
)

type Registration struct {
	gorm.Model
	CompetitionID uint `gorm:"not null;index"`
	Competition   Competition
	ParticipantID uint      `gorm:"not null"`
	Participant   User      `gorm:"foreignKey:ParticipantID"`
	PaymentStatus string    `gorm:"type:varchar(20)"` // pending | paid
	RegisteredAt  time.Time `gorm:"autoCreateTime"`
}
