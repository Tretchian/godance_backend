package competition

import (
	"strings"
	"time"

	"godance/internal/domain"
	"godance/internal/dto"
	"godance/internal/models"
	types "godance/internal/type"
)

type Service struct {
	repo domain.CompetitionRepository
}

func NewService(repo domain.CompetitionRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCompetition(organizerID uint, req dto.CreateCompetitionRequest) (*dto.Competition, error) {
	fee, _ := req.EntryFee.Float64()
	competition := models.Competition{
		OrganizerID:      organizerID,
		Title:            req.Title,
		EventDate:        req.EventDate.Time,
		ParticipantLimit: req.ParticipantLimit,
		EntryFee:         &fee,
		Status:           string(types.CompetitionStatusDraft),
		PayoutStatus:     string(types.PayoutStatusPending),
	}
	if err := s.repo.Create(&competition); err != nil {
		return nil, err
	}

	full, err := s.repo.GetByID(competition.ID)
	if err != nil {
		return nil, err
	}
	return toCompetition(full), nil
}

func (s *Service) GetCompetition(id uint) (*dto.Competition, error) {
	competition, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return toCompetition(competition), nil
}

func (s *Service) PatchCompetition(id, callerID uint, req dto.UpdateCompetitionRequest) (*dto.Competition, error) {
	competition, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if competition.OrganizerID != callerID {
		return nil, domain.ErrNotOwner
	}

	fields := map[string]any{}
	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.ParticipantLimit != nil {
		fields["participant_limit"] = *req.ParticipantLimit
	}
	if req.Status != nil {
		fields["status"] = string(*req.Status)
	}
	if len(fields) > 0 {
		if err := s.repo.UpdateFields(id, fields); err != nil {
			return nil, err
		}
	}

	updated, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return toCompetition(updated), nil
}

func (s *Service) Summary(id, callerID uint) (*dto.CompetitionSummary, error) {
	competition, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if competition.OrganizerID != callerID {
		return nil, domain.ErrNotOwner
	}

	counts, err := s.repo.FeedbackStatusCounts(id)
	if err != nil {
		return nil, err
	}
	captured, err := s.repo.CapturedAmount(id)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, c := range counts {
		total += c
	}
	completed := counts[string(types.FeedbackRequestStatusCompleted)]
	refunded := counts[string(types.FeedbackRequestStatusRefunded)]
	pending := total - completed - refunded

	var share float64
	if competition.FeedbackFeePercent != nil {
		share = captured * (*competition.FeedbackFeePercent) / 100
	}

	return &dto.CompetitionSummary{
		CompetitionID:       id,
		TotalRequests:       total,
		Completed:           completed,
		Pending:             pending,
		Refunded:            refunded,
		TotalCapturedAmount: captured,
		OrganizerShare:      share,
		PayoutStatus:        competition.PayoutStatus,
	}, nil
}

func (s *Service) GetPage(status string, page int, limit int) (*dto.CompetitionListResponse, error) {
	competitions, total, err := s.repo.FindPage(status, page, limit)
	if err != nil {
		return nil, err
	}

	items := make([]dto.Competition, len(competitions))
	for i := range competitions {
		items[i] = *toCompetition(&competitions[i])
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

func toCompetition(c *models.Competition) *dto.Competition {
	name, surname := splitName(organizerFullName(c))
	return &dto.Competition{
		ID:               c.ID,
		OrganizerID:      c.OrganizerID,
		OrganizerName:    name,
		OrganizerSurname: surname,
		Title:            c.Title,
		EventDate:        c.EventDate.Format("2006-01-02"),
		ParticipantLimit: c.ParticipantLimit,
		EntryFee:         c.EntryFee,
		Status:           c.Status,
		PayoutStatus:     c.PayoutStatus,
		CreatedAt:        c.CreatedAt.Format(time.RFC3339),
	}
}

func organizerFullName(c *models.Competition) string {
	if c.Organizer.Profile != nil {
		return c.Organizer.Profile.FullName
	}
	return ""
}

// splitName делит ФИО на имя и остальное (фамилию), т.к. в профиле хранится
// единое full_name, а схема Competition разделяет name/surname.
func splitName(fullName string) (name, surname string) {
	parts := strings.SplitN(strings.TrimSpace(fullName), " ", 2)
	name = parts[0]
	if len(parts) > 1 {
		surname = parts[1]
	}
	return name, surname
}
