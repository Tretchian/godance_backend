package user

import (
	"errors"

	"godance/internal/domain"
	"godance/internal/models"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.UserRepository {
	return &repository{db: db}
}

func (r *repository) Create(c *models.User) error {
	return r.db.Create(c).Error
}

func (r *repository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *repository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Preload("Profile").First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
