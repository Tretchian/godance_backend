package feedback

import (
	"errors"
	"os"
	"strconv"
	"time"

	"godance/internal/domain"
	"godance/internal/dto"
	"godance/internal/models"
	types "godance/internal/type"
)

const (
	defaultFeedbackPrice = 500.0
	requestTTL           = 30 * 24 * time.Hour // deadline_at = created_at + 30 дней
)

type Service struct {
	repo    domain.FeedbackRepository
	gateway domain.PaymentGateway
}

func NewService(repo domain.FeedbackRepository, gateway domain.PaymentGateway) *Service {
	return &Service{repo: repo, gateway: gateway}
}

// CreateRequest создаёт запрос ОС и инициирует hold оплаты через шлюз.
func (s *Service) CreateRequest(participantID uint, req dto.CreateFeedbackRequest) (*dto.FeedbackRequestCreated, error) {
	price := feedbackPrice()
	deadline := time.Now().Add(requestTTL)

	created, err := s.repo.CreateRequest(req.CompetitionID, req.JudgeID, participantID, req.Comment, price, deadline)
	if err != nil {
		return nil, err
	}

	paymentURL, err := s.gateway.CreateHold(created.ID, price)
	if err != nil {
		return nil, err
	}

	return &dto.FeedbackRequestCreated{
		RequestID:  created.ID,
		PaymentURL: paymentURL,
		Status:     created.Status,
	}, nil
}

// GetRequest доступен участнику-владельцу или судье-исполнителю запроса.
func (s *Service) GetRequest(id, callerID uint) (*dto.FeedbackRequestItem, error) {
	request, err := s.repo.GetRequest(id)
	if err != nil {
		return nil, err
	}
	if request.ParticipantID != callerID && request.JudgeID != callerID {
		return nil, domain.ErrNotRequestParticipant
	}
	return toRequestItem(request), nil
}

// ConfirmRequest — участник подтверждает полученную ОС, инициируя capture судье.
func (s *Service) ConfirmRequest(id, participantID uint) (*dto.FeedbackRequestItem, error) {
	request, err := s.repo.GetRequest(id)
	if err != nil {
		return nil, err
	}
	if request.ParticipantID != participantID {
		return nil, domain.ErrNotRequestParticipant
	}

	updated, err := s.repo.Transition(
		id,
		[]string{string(types.FeedbackRequestStatusAwaitingConfirmation)},
		string(types.FeedbackRequestStatusCompleted),
	)
	if err != nil {
		return nil, err
	}

	if err := s.gateway.Capture(id); err != nil {
		return nil, err
	}
	return toRequestItem(updated), nil
}

// CreateResponse — судья-исполнитель отправляет ОС (pending → awaiting_confirmation).
func (s *Service) CreateResponse(judgeID uint, req dto.CreateFeedbackResponse) (*dto.FeedbackResponseItem, error) {
	request, err := s.repo.GetRequest(req.RequestID)
	if err != nil {
		return nil, err
	}
	if request.JudgeID != judgeID {
		return nil, domain.ErrNotRequestJudge
	}

	response := models.FeedbackResponse{
		RequestID:       req.RequestID,
		Strengths:       req.Strengths,
		Errors:          req.Errors,
		Recommendations: req.Recommendations,
	}
	if err := s.repo.CreateResponse(&response); err != nil {
		return nil, err
	}
	return toResponseItem(&response), nil
}

// GetResponse доступен участнику-владельцу или судье-автору.
func (s *Service) GetResponse(id, callerID uint) (*dto.FeedbackResponseItem, error) {
	response, err := s.repo.GetResponse(id)
	if err != nil {
		return nil, err
	}
	request, err := s.repo.GetRequest(response.RequestID)
	if err != nil {
		return nil, err
	}
	if request.ParticipantID != callerID && request.JudgeID != callerID {
		return nil, domain.ErrNotRequestParticipant
	}
	return toResponseItem(response), nil
}

// CreateRating — участник оценивает ОС после подтверждения; рейтинг судьи пересчитывается.
func (s *Service) CreateRating(participantID uint, req dto.CreateFeedbackRating) (*dto.FeedbackRatingItem, error) {
	rating := models.FeedbackRating{
		ResponseID:    req.ResponseID,
		ParticipantID: participantID,
		Score:         int16(req.Score),
	}
	newRating, err := s.repo.CreateRating(&rating)
	if err != nil {
		return nil, err
	}
	return &dto.FeedbackRatingItem{
		ID:             rating.ID,
		ResponseID:     rating.ResponseID,
		Score:          req.Score,
		JudgeNewRating: newRating,
		RatedAt:        rating.RatedAt.Format(time.RFC3339),
	}, nil
}

// HandleFeedbackWebhook драйвит статусы запроса по событиям платёжного шлюза.
// Идемпотентен: повторные/устаревшие события не считаются ошибкой.
func (s *Service) HandleFeedbackWebhook(p dto.PaymentWebhook) error {
	var allowedFrom []string
	var to string

	switch types.PaymentWebhookStatus(p.Status) {
	case types.PaymentWebhookStatusHeld:
		allowedFrom = []string{string(types.FeedbackRequestStatusAwaitingPayment)}
		to = string(types.FeedbackRequestStatusAwaitingVideo)
	case types.PaymentWebhookStatusFailed, types.PaymentWebhookStatusCaptureFailed:
		allowedFrom = []string{string(types.FeedbackRequestStatusAwaitingPayment)}
		to = string(types.FeedbackRequestStatusRefunded)
	case types.PaymentWebhookStatusReleased:
		allowedFrom = []string{
			string(types.FeedbackRequestStatusAwaitingVideo),
			string(types.FeedbackRequestStatusPending),
			string(types.FeedbackRequestStatusAwaitingConfirmation),
		}
		to = string(types.FeedbackRequestStatusRefunded)
	case types.PaymentWebhookStatusCaptured:
		allowedFrom = []string{
			string(types.FeedbackRequestStatusAwaitingConfirmation),
			string(types.FeedbackRequestStatusCompleted),
		}
		to = string(types.FeedbackRequestStatusCompleted)
	default:
		return nil // succeeded / paid_out и пр. — для feedback не несут перехода
	}

	_, err := s.repo.Transition(p.RelatedID, allowedFrom, to)
	if errors.Is(err, domain.ErrInvalidStatusTransition) {
		return nil // идемпотентность: событие уже учтено ранее
	}
	return err
}

func toRequestItem(r *models.FeedbackRequest) *dto.FeedbackRequestItem {
	item := &dto.FeedbackRequestItem{
		ID:            r.ID,
		JudgeID:       r.JudgeID,
		ParticipantID: r.ParticipantID,
		VideoID:       r.VideoID,
		Comment:       r.ParticipantComment,
		Price:         r.Price,
		Status:        r.Status,
		CreatedAt:     r.CreatedAt.Format(time.RFC3339),
	}
	if r.DeadlineAt != nil {
		deadline := r.DeadlineAt.Format(time.RFC3339)
		item.DeadlineAt = &deadline
	}
	return item
}

func toResponseItem(r *models.FeedbackResponse) *dto.FeedbackResponseItem {
	return &dto.FeedbackResponseItem{
		ID:              r.ID,
		RequestID:       r.RequestID,
		Strengths:       r.Strengths,
		Errors:          r.Errors,
		Recommendations: r.Recommendations,
		SubmittedAt:     r.CreatedAt.Format(time.RFC3339),
	}
}

func feedbackPrice() float64 {
	if v := os.Getenv("FEEDBACK_PRICE"); v != "" {
		if p, err := strconv.ParseFloat(v, 64); err == nil {
			return p
		}
	}
	return defaultFeedbackPrice
}
