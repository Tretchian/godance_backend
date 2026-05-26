package user

import (
	"godance/internal/domain"
)

type Service struct {
	repo domain.UserRepository
}

func NewService(repo domain.UserRepository) *Service {
	return &Service{repo: repo}
}
