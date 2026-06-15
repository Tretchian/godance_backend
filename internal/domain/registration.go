package domain

import "godance/internal/models"

type RegistrationRepository interface {
	ListByCompetition(competitionID, callerID uint, page, limit int) ([]models.Registration, int64, error)
	Register(competitionID, participantID uint) (*models.Registration, error)
}
