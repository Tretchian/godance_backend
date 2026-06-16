package storage

import (
	"context"
	"time"

	"godance/internal/domain"
)

// Stub — no-op реализация VideoStorage для тестов и локального запуска без
// поднятого MinIO. Возвращает фиктивные URL и считает любой объект существующим.
type Stub struct{}

func NewStub() *Stub { return &Stub{} }

var _ domain.VideoStorage = (*Stub)(nil)

func (s *Stub) PresignPut(_ context.Context, key, _ string, _ time.Duration) (string, error) {
	return "https://stub.local/put/" + key, nil
}

func (s *Stub) PresignGet(_ context.Context, key string, _ time.Duration) (string, error) {
	return "https://stub.local/get/" + key, nil
}

func (s *Stub) Stat(_ context.Context, _ string) (int64, error) {
	return 1, nil
}
