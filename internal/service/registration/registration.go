package registration

import (
	"fmt"
	"os"

	"godance/internal/domain"
	"godance/internal/dto"
	"godance/internal/models"
)

type Service struct {
	repo domain.RegistrationRepository
}

func NewService(repo domain.RegistrationRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListForCompetition(competitionID, callerID uint, page, limit int) (*dto.RegistrationListResponse, error) {
	registrations, total, err := s.repo.ListByCompetition(competitionID, callerID, page, limit)
	if err != nil {
		return nil, err
	}

	items := make([]dto.RegistrationItem, len(registrations))
	for i, reg := range registrations {
		items[i] = toRegistrationItem(reg)
	}

	return &dto.RegistrationListResponse{
		Data: items,
		Pagination: dto.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (s *Service) Register(competitionID, participantID uint) (*dto.RegistrationCreated, error) {
	registration, err := s.repo.Register(competitionID, participantID)
	if err != nil {
		return nil, err
	}

	return &dto.RegistrationCreated{
		RegistrationID: registration.ID,
		PaymentURL:     paymentURL(registration.ID),
	}, nil
}

func toRegistrationItem(r models.Registration) dto.RegistrationItem {
	fullName := ""
	if r.Participant.Profile != nil {
		fullName = r.Participant.Profile.FullName
	}
	return dto.RegistrationItem{
		ID:            r.ID,
		ParticipantID: r.ParticipantID,
		FullName:      fullName,
		PaymentStatus: r.PaymentStatus,
		RegisteredAt:  r.RegisteredAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// paymentURL — временная заглушка ссылки на платёжный шлюз.
// Будет заменена реальной интеграцией в срезе payments.
func paymentURL(registrationID uint) string {
	return fmt.Sprintf("%s/pay/registration/%d", os.Getenv("PAYMENT_GATEWAY_URL"), registrationID)
}
