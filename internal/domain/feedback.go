package domain

import (
	"time"

	"godance/internal/models"
)

type FeedbackRepository interface {
	// CreateRequest создаёт запрос ОС, предварительно проверяя, что
	// соревнование существует и судья на него назначен.
	CreateRequest(competitionID, judgeID, participantID uint, comment string, price float64, deadline time.Time) (*models.FeedbackRequest, error)
	GetRequest(id uint) (*models.FeedbackRequest, error)
	// Transition атомарно меняет статус запроса под блокировкой строки,
	// допуская переход только из одного из allowedFrom состояний.
	Transition(id uint, allowedFrom []string, to string) (*models.FeedbackRequest, error)
	// AttachVideo привязывает видео к запросу участника-владельца и переводит
	// его awaiting_video → pending в одной транзакции; comment, если задан,
	// обновляет комментарий участника.
	RequestByVideo(videoID uint) (*models.FeedbackRequest, error)
	AttachVideo(requestID, videoID, participantID uint, comment *string) error

	// CreateResponse сохраняет ОС судьи и переводит запрос pending →
	// awaiting_confirmation в одной транзакции.
	CreateResponse(resp *models.FeedbackResponse) error
	GetResponse(id uint) (*models.FeedbackResponse, error)

	// CreateRating сохраняет оценку и пересчитывает средний рейтинг судьи,
	// возвращая обновлённое значение.
	CreateRating(rating *models.FeedbackRating) (judgeNewRating float64, err error)
}
