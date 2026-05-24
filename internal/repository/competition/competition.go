package competition

import (
	"godance/internal/domain"
	"godance/internal/models"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.CompetitionRepository {
	return &repository{db: db}
}

func (r *repository) Create(c *models.Competition) error {
	return r.db.Create(c).Error
}

func (r *repository) FindPage(status string, page int, limit int) ([]models.Competition, int64, error) {
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
