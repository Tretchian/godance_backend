package domain

import (
	"context"
	"time"

	"godance/internal/models"
)

type VideoRepository interface {
	Create(v *models.Video) error
	Get(id uint) (*models.Video, error)
	// Activate переводит видео uploading → active и проставляет срок хранения.
	Activate(id uint, expiresAt time.Time) (*models.Video, error)
}

// VideoStorage абстрагирует S3-совместимое объектное хранилище (MinIO).
type VideoStorage interface {
	// PresignPut выдаёт временный URL для прямой загрузки клиентом (PUT).
	PresignPut(ctx context.Context, key, contentType string, ttl time.Duration) (string, error)
	// PresignGet выдаёт временный URL для просмотра (GET).
	PresignGet(ctx context.Context, key string, ttl time.Duration) (string, error)
	// Stat проверяет наличие объекта (HEAD) и возвращает его размер.
	// Возвращает ErrVideoObjectMissing, если объекта нет.
	Stat(ctx context.Context, key string) (size int64, err error)
}
