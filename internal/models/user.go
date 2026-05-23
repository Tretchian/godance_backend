package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email        string `gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Phone        string `gorm:"type:varchar(20);uniqueIndex;not null"`
	Role         string `gorm:"type:varchar(20);not null"` // participant | judge | organizer

	Profile       *Profile
	Competitions  []Competition  `gorm:"foreignKey:OrganizerID"`
	Registrations []Registration `gorm:"foreignKey:ParticipantID"`
	Videos        []Video        `gorm:"foreignKey:ParticipantID"`
	Notifications []Notification
	Payments      []Payment
	AuditLogs     []AuditLog
}
