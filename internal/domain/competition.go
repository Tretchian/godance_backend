package domain

import (
	"godance/internal/dto"
	"godance/internal/models"
)

type CompetitionRepository interface {
	Create(c *models.Competition) error
	FindPage(status string, page int, limit int) ([]models.Competition, int64, error)
	GetByID(id int) (*models.Competition, error)
	PatchByID(id int, c *dto.UpdateCompetitionRequest) (*models.Competition, error)
	GetSummaryById(id int) (dto.CompetitionSummary, error)
}
