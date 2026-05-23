package feedback

import "godance/internal/repository/feedback"

type Service struct {
	repo *feedback.Repository
}

func NewService(repo *feedback.Repository) *Service {
	return &Service{repo: repo}
}
