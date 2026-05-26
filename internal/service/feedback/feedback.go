package feedback

import (
	"godance/internal/domain"
)

type Service struct {
	repo domain.FeedbackRepository
}

func NewService(repo domain.FeedbackRepository) *Service {
	return &Service{repo: repo}
}
