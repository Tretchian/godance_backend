package competition

import (
	"godance/internal/domain"
	"godance/internal/dto"
	"godance/internal/models"

	"github.com/google/uuid"
)

type Service struct {
	repo domain.CompetitionRepository
}

func NewService(repo domain.CompetitionRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCompetition(user models.User, req dto.CreateCompetitionRequest) (*dto.CompetitionItem, error) {
	competition := models.Competition{
		OrganizerID: user.ID,
		Title:       req.Title,
		EventDate:   req.EventDate,
		Status:      "draft",
	}

	err := s.repo.Create(&competition)
	if err != nil {
		return nil, err
	}

	return toCompetitionItem(competition), nil
}

func (s *Service) PatchCompetition(id uuid.UUID, req dto.UpdateCompetitionRequest) (*dto.CompetitionItem, error) {
	return nil, nil
}

func (s *Service) GetPage(status string, page int, limit int) (*dto.CompetitionListResponse, error) {
	competitions, total, err := s.repo.FindPage(status, page, limit)
	if err != nil {
		return nil, err
	}

	items := make([]dto.CompetitionItem, len(competitions))

	for i, c := range competitions {
		items[i] = dto.CompetitionItem{
			OrganizerID:        c.OrganizerID,
			Title:              c.Title,
			EventDate:          c.EventDate.Format("2006-01-02"),
			ParticipantLimit:   c.ParticipantLimit,
			EntryFee:           c.EntryFee,
			FeedbackFeePercent: c.FeedbackFeePercent,
			Status:             c.Status,
		}
	}

	return &dto.CompetitionListResponse{
		Data: items,
		Pagination: dto.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func toCompetitionItem(c models.Competition) *dto.CompetitionItem {
	return &dto.CompetitionItem{
		Title:     c.Title,
		EventDate: c.EventDate.Format("2006-01-02"),
		Status:    c.Status,
	}
}
