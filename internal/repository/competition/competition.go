package competition

import (
	"godance/internal/domain"
	"godance/internal/dto"
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

func (r *repository) GetByID(id int) (*models.Competition, error) {
	var competition models.Competition
	err := r.db.First(&competition, id).Error
	if err != nil {
		return nil, err
	}
	return &competition, nil
}

func (r *repository) GetSummaryById(id int) (dto.CompetitionSummary, error) {
	var summary dto.CompetitionSummary

	return summary, nil
}

func (r *repository) PatchByID(id int, c *dto.UpdateCompetitionRequest) (*models.Competition, error) {
	var competition models.Competition

	if err := r.db.First(&competition, id).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&competition).Updates(c).Error; err != nil {
		return nil, err
	}

	return &competition, nil
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
