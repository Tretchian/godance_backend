package domain

import "godance/internal/models"

type RegistrationRepository interface {
	ListByCompetition(competitionID, callerID uint, page, limit int) ([]models.Registration, int64, error)
	Register(competitionID, participantID uint) (*models.Registration, error)
	// MarkPaid переводит payment_status регистрации в paid (вызывается из
	// вебхука платёжного шлюза).
	MarkPaid(id uint) error
}
