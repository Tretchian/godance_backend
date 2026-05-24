package payment

import (
	"godance/internal/domain"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.PaymentRepository {
	return &repository{db: db}
}
