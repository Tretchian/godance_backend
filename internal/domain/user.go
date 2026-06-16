package domain

import "godance/internal/models"

type UserRepository interface {
	Create(c *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
}
