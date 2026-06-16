package models

import (
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	UserID   uint `gorm:"uniqueIndex;not null"`
	User     User
	FullName string  `gorm:"type:varchar(255)"`
	Bio      string  `gorm:"type:text"`
	Rating   float64 `gorm:"type:numeric(3,2);default:0"`
}
