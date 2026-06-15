package models

import (
	"time"

	"gorm.io/gorm"
)

type Registration struct {
	gorm.Model
	CompetitionID uint `gorm:"not null;uniqueIndex:idx_comp_participant;index"`
	Competition   Competition
	ParticipantID uint      `gorm:"not null;uniqueIndex:idx_comp_participant"`
	Participant   User      `gorm:"foreignKey:ParticipantID"`
	PaymentStatus string    `gorm:"type:varchar(20)"` // pending | paid
	RegisteredAt  time.Time `gorm:"autoCreateTime"`
}
