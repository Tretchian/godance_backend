package domain

import "godance/internal/models"

type JudgeRepository interface {
	ListByCompetition(competitionID uint) ([]models.JudgesCompetition, error)
	Assign(competitionID, judgeID, callerID uint, limit int) (*models.JudgesCompetition, error)
}
