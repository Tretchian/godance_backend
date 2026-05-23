package user

import (
	"godance/internal/repository/user"
)

type Service struct {
	repo *user.Repository
}

func NewService(repo *user.Repository) *Service {
	return &Service{repo: repo}
}
