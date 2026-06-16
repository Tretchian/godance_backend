package registration

import (
	"errors"

	"godance/internal/domain"
	"godance/internal/models"
	types "godance/internal/type"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.RegistrationRepository {
	return &repository{db: db}
}

// ListByCompetition возвращает регистрации соревнования с пагинацией.
// Доступно только организатору-владельцу соревнования.
func (r *repository) ListByCompetition(competitionID, callerID uint, page, limit int) ([]models.Registration, int64, error) {
	var comp models.Competition
	if err := r.db.First(&comp, competitionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, domain.ErrCompetitionNotFound
		}
		return nil, 0, err
	}
	if comp.OrganizerID != callerID {
		return nil, 0, domain.ErrNotOwner
	}

	var total int64
	if err := r.db.Model(&models.Registration{}).
		Where("competition_id = ?", competitionID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var registrations []models.Registration
	err := r.db.
		Preload("Participant.Profile"). // нужен full_name из profiles
		Where("competition_id = ?", competitionID).
		Order("registered_at").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&registrations).Error
	return registrations, total, err
}

// Register — атомарная регистрация участника под блокировкой строки
// соревнования. Инвариант "не более participant_limit" удерживается за счёт
// сериализации параллельных регистраций через FOR UPDATE на строке-родителе.
func (r *repository) Register(competitionID, participantID uint) (*models.Registration, error) {
	var registration models.Registration

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var comp models.Competition
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&comp, competitionID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrCompetitionNotFound
			}
			return err
		}

		// Регистрироваться можно только на открытые соревнования.
		if comp.Status != string(types.CompetitionStatusOpen) {
			return domain.ErrRegistrationClosed
		}

		// Лимит участников под удерживаемой блокировкой родителя (nil = без лимита).
		if comp.ParticipantLimit != nil {
			var count int64
			if err := tx.Model(&models.Registration{}).
				Where("competition_id = ?", competitionID).
				Count(&count).Error; err != nil {
				return err
			}
			if count >= int64(*comp.ParticipantLimit) {
				return domain.ErrParticipantLimitReached
			}
		}

		registration = models.Registration{
			CompetitionID: competitionID,
			ParticipantID: participantID,
			PaymentStatus: string(types.RegistrationPaymentStatusPending),
		}
		if err := tx.Create(&registration).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return domain.ErrAlreadyRegistered
			}
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := r.db.Preload("Participant.Profile").
		First(&registration, registration.ID).Error; err != nil {
		return nil, err
	}
	return &registration, nil
}

func (r *repository) MarkPaid(id uint) error {
	res := r.db.Model(&models.Registration{}).
		Where("id = ?", id).
		Update("payment_status", string(types.RegistrationPaymentStatusPaid))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrRegistrationNotFound
	}
	return nil
}
