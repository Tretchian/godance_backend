package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	UserID  uint `gorm:"not null"`
	User    User
	Channel string    `gorm:"type:varchar(10)"` // sms
	Message *string   `gorm:"type:text"`
	Status  string    `gorm:"type:varchar(20)"` // sent | failed
	SentAt  time.Time `gorm:"autoCreateTime"`
}
