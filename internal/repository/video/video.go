package video

import (
	"errors"
	"time"

	"godance/internal/domain"
	"godance/internal/models"
	types "godance/internal/type"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.VideoRepository {
	return &repository{db: db}
}

func (r *repository) Create(v *models.Video) error {
	return r.db.Create(v).Error
}

func (r *repository) Get(id uint) (*models.Video, error) {
	var video models.Video
	if err := r.db.First(&video, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrVideoNotFound
		}
		return nil, err
	}
	return &video, nil
}

func (r *repository) Activate(id uint, expiresAt time.Time) (*models.Video, error) {
	var video models.Video
	if err := r.db.First(&video, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrVideoNotFound
		}
		return nil, err
	}
	video.Status = string(types.VideoStatusActive)
	video.ExpiresAt = expiresAt
	if err := r.db.Model(&video).Updates(map[string]any{
		"status":     video.Status,
		"expires_at": video.ExpiresAt,
	}).Error; err != nil {
		return nil, err
	}
	return &video, nil
}
