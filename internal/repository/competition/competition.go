package competition

import (
	"godance/internal/models"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(c *models.Competition) error {
	return r.db.Create(c).Error
}

func (r *Repository) FindPage(status string, page int, limit int) ([]models.Competition, int64, error) {
	var competitions []models.Competition
	var total int64

	query := r.db.Model(&models.Competition{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset((page - 1) * limit).Limit(limit).Find(&competitions).Error

	return competitions, total, err
}
