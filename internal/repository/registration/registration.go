package registration

import (
	"godance/internal/domain"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.RegistrationRepository {
	return &repository{db: db}
}
