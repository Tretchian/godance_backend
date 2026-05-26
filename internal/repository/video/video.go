package video

import (
	"godance/internal/domain"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.VideoRepository {
	return &repository{db: db}
}
