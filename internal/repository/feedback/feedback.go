package feedback

import (
	"errors"
	"slices"
	"time"

	"godance/internal/domain"
	"godance/internal/models"
	types "godance/internal/type"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.FeedbackRepository {
	return &repository{db: db}
}

func (r *repository) CreateRequest(competitionID, judgeID, participantID uint, comment string, price float64, deadline time.Time) (*models.FeedbackRequest, error) {
	request := models.FeedbackRequest{
		CompetitionID:      competitionID,
		JudgeID:            judgeID,
		ParticipantID:      participantID,
		Price:              &price,
		ParticipantComment: &comment,
		Status:             string(types.FeedbackRequestStatusAwaitingPayment),
		DeadlineAt:         &deadline,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var comp models.Competition
		if err := tx.First(&comp, competitionID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrCompetitionNotFound
			}
			return err
		}

		var assigned int64
		if err := tx.Model(&models.JudgesCompetition{}).
			Where("competition_id = ? AND judge_id = ?", competitionID, judgeID).
			Count(&assigned).Error; err != nil {
			return err
		}
		if assigned == 0 {
			return domain.ErrJudgeNotAssigned
		}

		return tx.Create(&request).Error
	})
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *repository) GetRequest(id uint) (*models.FeedbackRequest, error) {
	var request models.FeedbackRequest
	if err := r.db.First(&request, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrFeedbackRequestNotFound
		}
		return nil, err
	}
	return &request, nil
}

func (r *repository) ListExpired(now time.Time) ([]models.FeedbackRequest, error) {
	active := []string{
		string(types.FeedbackRequestStatusAwaitingPayment),
		string(types.FeedbackRequestStatusAwaitingVideo),
		string(types.FeedbackRequestStatusPending),
		string(types.FeedbackRequestStatusAwaitingConfirmation),
	}
	var requests []models.FeedbackRequest
	err := r.db.
		Where("deadline_at IS NOT NULL AND deadline_at < ? AND status IN ?", now, active).
		Find(&requests).Error
	return requests, err
}

func (r *repository) Transition(id uint, allowedFrom []string, to string) (*models.FeedbackRequest, error) {
	var request models.FeedbackRequest
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&request, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrFeedbackRequestNotFound
			}
			return err
		}
		if !slices.Contains(allowedFrom, request.Status) {
			return domain.ErrInvalidStatusTransition
		}
		request.Status = to
		return tx.Model(&request).Update("status", to).Error
	})
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *repository) RequestByVideo(videoID uint) (*models.FeedbackRequest, error) {
	var request models.FeedbackRequest
	if err := r.db.Where("video_id = ?", videoID).First(&request).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrFeedbackRequestNotFound
		}
		return nil, err
	}
	return &request, nil
}

func (r *repository) AttachVideo(requestID, videoID, participantID uint, comment *string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var request models.FeedbackRequest
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&request, requestID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrFeedbackRequestNotFound
			}
			return err
		}
		if request.ParticipantID != participantID {
			return domain.ErrNotRequestParticipant
		}
		if request.Status != string(types.FeedbackRequestStatusAwaitingVideo) {
			return domain.ErrRequestNotAwaitingVideo
		}

		updates := map[string]any{
			"video_id": videoID,
			"status":   string(types.FeedbackRequestStatusPending),
		}
		if comment != nil {
			updates["participant_comment"] = *comment
		}
		return tx.Model(&request).Updates(updates).Error
	})
}

func (r *repository) CreateResponse(resp *models.FeedbackResponse) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var request models.FeedbackRequest
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&request, resp.RequestID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrFeedbackRequestNotFound
			}
			return err
		}
		if request.Status != string(types.FeedbackRequestStatusPending) {
			return domain.ErrInvalidStatusTransition
		}

		if err := tx.Create(resp).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return domain.ErrResponseAlreadyExists
			}
			return err
		}

		return tx.Model(&request).
			Update("status", string(types.FeedbackRequestStatusAwaitingConfirmation)).Error
	})
}

func (r *repository) GetResponse(id uint) (*models.FeedbackResponse, error) {
	var response models.FeedbackResponse
	if err := r.db.First(&response, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrFeedbackResponseNotFound
		}
		return nil, err
	}
	return &response, nil
}

func (r *repository) CreateRating(rating *models.FeedbackRating) (float64, error) {
	var judgeID uint
	var newRating float64

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var response models.FeedbackResponse
		if err := tx.First(&response, rating.ResponseID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrFeedbackResponseNotFound
			}
			return err
		}

		var request models.FeedbackRequest
		if err := tx.First(&request, response.RequestID).Error; err != nil {
			return err
		}
		// Оценить можно только подтверждённую (завершённую) ОС и только её владельцу.
		if request.Status != string(types.FeedbackRequestStatusCompleted) {
			return domain.ErrInvalidStatusTransition
		}
		if request.ParticipantID != rating.ParticipantID {
			return domain.ErrNotRequestParticipant
		}
		judgeID = request.JudgeID

		if err := tx.Create(rating).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return domain.ErrRatingAlreadyExists
			}
			return err
		}

		// Пересчёт среднего рейтинга судьи по всем его оценённым ОС.
		if err := tx.Model(&models.FeedbackRating{}).
			Joins("JOIN feedback_responses resp ON resp.id = feedback_ratings.response_id").
			Joins("JOIN feedback_requests req ON req.id = resp.request_id").
			Where("req.judge_id = ?", judgeID).
			Select("COALESCE(AVG(feedback_ratings.score), 0)").
			Scan(&newRating).Error; err != nil {
			return err
		}

		return tx.Model(&models.Profile{}).
			Where("user_id = ?", judgeID).
			Update("rating", newRating).Error
	})
	if err != nil {
		return 0, err
	}
	return newRating, nil
}
