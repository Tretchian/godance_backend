package competition

import (
	"errors"

	"godance/internal/domain"
	"godance/internal/models"
	types "godance/internal/type"

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

	err := query.
		Preload("Organizer.Profile").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&competitions).Error

	return competitions, total, err
}

func (r *repository) GetByID(id uint) (*models.Competition, error) {
	var competition models.Competition
	if err := r.db.Preload("Organizer.Profile").First(&competition, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrCompetitionNotFound
		}
		return nil, err
	}
	return &competition, nil
}

func (r *repository) UpdateFields(id uint, fields map[string]any) error {
	return r.db.Model(&models.Competition{}).Where("id = ?", id).Updates(fields).Error
}

func (r *repository) FeedbackStatusCounts(competitionID uint) (map[string]int64, error) {
	type row struct {
		Status string
		Cnt    int64
	}
	var rows []row
	if err := r.db.Model(&models.FeedbackRequest{}).
		Select("status, count(*) as cnt").
		Where("competition_id = ?", competitionID).
		Group("status").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	counts := make(map[string]int64, len(rows))
	for _, x := range rows {
		counts[x.Status] = x.Cnt
	}
	return counts, nil
}

func (r *repository) CapturedAmount(competitionID uint) (float64, error) {
	var sum float64
	err := r.db.Model(&models.FeedbackRequest{}).
		Where("competition_id = ? AND status = ?", competitionID, string(types.FeedbackRequestStatusCompleted)).
		Select("COALESCE(SUM(price), 0)").
		Scan(&sum).Error
	return sum, err
}
