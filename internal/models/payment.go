package models

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	UserID      uint `gorm:"not null"`
	User        User
	RelatedID   uint      `gorm:"not null"`         // ID регистрации или запроса ОС
	RelatedType *string   `gorm:"type:varchar(20)"` // registration | feedback | organizer_payout
	Amount      float64   `gorm:"type:numeric(10,2)"`
	Status      string    `gorm:"type:varchar(20)"` // pending | succeeded | refunded | failed
	GatewayRef  string    `gorm:"type:varchar"`     // ID транзакции у шлюза
	PaidAt      time.Time `gorm:"type:timestamptz"`
}
