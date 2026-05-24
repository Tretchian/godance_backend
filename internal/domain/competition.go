package domain

import "godance/internal/models"

type CompetitionRepository interface {
	Create(c *models.Competition) error
	FindPage(status string, page int, limit int) ([]models.Competition, int64, error)
}
