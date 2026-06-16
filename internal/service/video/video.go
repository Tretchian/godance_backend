package video

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"godance/internal/domain"
	"godance/internal/dto"
	"godance/internal/models"
	types "godance/internal/type"
)

const (
	uploadTTL = time.Hour           // presigned PUT
	viewTTL   = 15 * time.Minute    // presigned GET
	retention = 30 * 24 * time.Hour // expires_at = now + 1 месяц
)

type Service struct {
	repo     domain.VideoRepository
	storage  domain.VideoStorage
	feedback domain.FeedbackRepository
}

func NewService(repo domain.VideoRepository, storage domain.VideoStorage, feedback domain.FeedbackRepository) *Service {
	return &Service{repo: repo, storage: storage, feedback: feedback}
}

// CreateUploadURL выдаёт presigned PUT URL. Доступно только владельцу запроса
// в статусе awaiting_video.
func (s *Service) CreateUploadURL(participantID uint, req dto.UploadUrlRequest) (*dto.UploadUrlResponse, error) {
	request, err := s.feedback.GetRequest(req.RequestID)
	if err != nil {
		return nil, err
	}
	if request.ParticipantID != participantID {
		return nil, domain.ErrNotRequestParticipant
	}
	if request.Status != string(types.FeedbackRequestStatusAwaitingVideo) {
		return nil, domain.ErrRequestNotAwaitingVideo
	}

	key, err := storageKey(participantID, req.ContentType)
	if err != nil {
		return nil, err
	}

	video := models.Video{
		ParticipantID: participantID,
		StorageKey:    &key,
		Status:        string(types.VideoStatusUploading),
		ContentType:   req.ContentType,
		SizeBytes:     req.FileSize,
	}
	if err := s.repo.Create(&video); err != nil {
		return nil, err
	}

	uploadURL, err := s.storage.PresignPut(context.Background(), key, req.ContentType, uploadTTL)
	if err != nil {
		return nil, err
	}

	return &dto.UploadUrlResponse{
		VideoID:   video.ID,
		UploadURL: uploadURL,
		ExpiresAt: time.Now().Add(uploadTTL).Format(time.RFC3339),
	}, nil
}

// ConfirmUpload проверяет наличие объекта в хранилище (HEAD), активирует видео
// и переводит запрос awaiting_video → pending.
func (s *Service) ConfirmUpload(participantID, videoID uint, req dto.VideoConfirmRequest) (*dto.VideoConfirmResponse, error) {
	video, err := s.repo.Get(videoID)
	if err != nil {
		return nil, err
	}
	if video.ParticipantID != participantID {
		return nil, domain.ErrNotVideoOwner
	}
	if video.StorageKey == nil {
		return nil, domain.ErrVideoObjectMissing
	}

	if _, err := s.storage.Stat(context.Background(), *video.StorageKey); err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(retention)
	activated, err := s.repo.Activate(videoID, expiresAt)
	if err != nil {
		return nil, err
	}

	if err := s.feedback.AttachVideo(req.RequestID, videoID, participantID, req.Comment); err != nil {
		return nil, err
	}

	return &dto.VideoConfirmResponse{
		VideoID:       activated.ID,
		Status:        activated.Status,
		ExpiresAt:     expiresAt.Format(time.RFC3339),
		RequestStatus: string(types.FeedbackRequestStatusPending),
	}, nil
}

// GetViewURL выдаёт presigned GET URL. Доступно владельцу-участнику или судье
// по запросу ОС, связанному с этим видео.
func (s *Service) GetViewURL(callerID, videoID uint) (*dto.VideoViewResponse, error) {
	video, err := s.repo.Get(videoID)
	if err != nil {
		return nil, err
	}
	if video.Status == string(types.VideoStatusDeleted) || video.StorageKey == nil {
		return nil, domain.ErrVideoGone
	}
	if !video.ExpiresAt.IsZero() && video.ExpiresAt.Before(time.Now()) {
		return nil, domain.ErrVideoGone
	}

	if video.ParticipantID != callerID {
		request, err := s.feedback.RequestByVideo(videoID)
		if err != nil || request.JudgeID != callerID {
			return nil, domain.ErrNotVideoOwner
		}
	}

	viewURL, err := s.storage.PresignGet(context.Background(), *video.StorageKey, viewTTL)
	if err != nil {
		return nil, err
	}

	return &dto.VideoViewResponse{
		VideoID:      video.ID,
		ViewURL:      viewURL,
		URLExpiresAt: time.Now().Add(viewTTL).Format(time.RFC3339),
		VideoLength:  "", // длительность не вычисляется (нет транскодинга)
	}, nil
}

func storageKey(participantID uint, contentType string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return fmt.Sprintf("videos/%d/%s%s", participantID, hex.EncodeToString(buf), ext(contentType)), nil
}

func ext(contentType string) string {
	switch contentType {
	case "video/quicktime":
		return ".mov"
	default:
		return ".mp4"
	}
}
