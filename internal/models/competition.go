package models

import (
	"time"

	"gorm.io/gorm"
)

type Competition struct {
	gorm.Model
	OrganizerID        uint      `gorm:"not null"`
	Organizer          User      `gorm:"foreignKey:OrganizerID"`
	Title              string    `gorm:"type:varchar(255);not null"`
	EventDate          time.Time `gorm:"type:date;not null"`
	ParticipantLimit   *int      `gorm:"type:int"`
	EntryFee           *float64  `gorm:"type:numeric(10,2)"`
	Status             string    `gorm:"type:varchar(20)"` // draft | open | closed | finished
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	FeedbackFeePercent *float64  `gorm:"type:numeric(5,2)"`
	PayoutStatus       string    `gorm:"type:varchar(20)"` // pending | completed | failed

	Registrations []Registration
}
