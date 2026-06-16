package domain

import "godance/internal/models"

type CompetitionRepository interface {
	Create(c *models.Competition) error
	FindPage(status string, page int, limit int) ([]models.Competition, int64, error)
	GetByID(id uint) (*models.Competition, error)
	UpdateFields(id uint, fields map[string]any) error
	// FeedbackStatusCounts возвращает количество запросов ОС соревнования,
	// сгруппированное по статусу.
	FeedbackStatusCounts(competitionID uint) (map[string]int64, error)
	// CapturedAmount возвращает сумму захваченных (completed) платежей по ОС.
	CapturedAmount(competitionID uint) (float64, error)
}
