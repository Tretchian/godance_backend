package judge

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

func NewRepository(db *gorm.DB) domain.JudgeRepository {
	return &repository{db: db}
}

func (r *repository) ListByCompetition(competitionID uint) ([]models.JudgesCompetition, error) {
	var assignments []models.JudgesCompetition
	err := r.db.
		Preload("Judge.Profile"). // нужны full_name и rating из profiles
		Where("competition_id = ?", competitionID).
		Order("assigned_at").
		Find(&assignments).Error
	return assignments, err
}

// Assign — атомарное назначение под блокировкой строки соревнования.
// Инвариант "не более limit судей" удерживается за счёт сериализации
// параллельных назначений через FOR UPDATE на строке-родителе: COUNT видит
// уже зафиксированные строки, а конкурирующая транзакция ждёт на этой же строке.
func (r *repository) Assign(competitionID, judgeID, callerID uint, limit int) (*models.JudgesCompetition, error) {
	var assignment models.JudgesCompetition

	err := r.db.Transaction(func(tx *gorm.DB) error {
		// Блокировка строки соревнования (SELECT ... FOR UPDATE).
		var comp models.Competition
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&comp, competitionID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrCompetitionNotFound
			}
			return err
		}

		// Можно назначать судей только на свои соревнования
		if comp.OrganizerID != callerID {
			return domain.ErrNotOwner
		}
		if comp.Status != string(types.CompetitionStatusDraft) &&
			comp.Status != string(types.CompetitionStatusOpen) {
			return domain.ErrCompetitionClosed
		}

		// Назначаемый обязан существовать и иметь роль judge
		var judge models.User
		if err := tx.First(&judge, judgeID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrUserNotFound
			}
			return err
		}
		if judge.Role != string(types.UserRoleJudge) {
			return domain.ErrNotJudge
		}

		// Лимит под удерживаемой блокировкой родителя
		var count int64
		if err := tx.Model(&models.JudgesCompetition{}).
			Where("competition_id = ?", competitionID).
			Count(&count).Error; err != nil {
			return err
		}
		if count >= int64(limit) {
			return domain.ErrJudgeLimitReached
		}

		// Проверка гонки данных в окне между COUNT и INSERT.
		assignment = models.JudgesCompetition{
			CompetitionID: competitionID,
			JudgeID:       judgeID,
		}
		if err := tx.Create(&assignment).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return domain.ErrAlreadyAssigned
			}
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := r.db.Preload("Judge.Profile").First(&assignment, assignment.ID).Error; err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (r *repository) Unassign(competitionID, judgeID, callerID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var comp models.Competition
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&comp, competitionID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrCompetitionNotFound
			}
			return err
		}
		if comp.OrganizerID != callerID {
			return domain.ErrNotOwner
		}

		// Снятие запрещено, если по судье уже есть запросы ОС.
		var requests int64
		if err := tx.Model(&models.FeedbackRequest{}).
			Where("competition_id = ? AND judge_id = ?", competitionID, judgeID).
			Count(&requests).Error; err != nil {
			return err
		}
		if requests > 0 {
			return domain.ErrJudgeHasRequests
		}

		res := tx.Where("competition_id = ? AND judge_id = ?", competitionID, judgeID).
			Delete(&models.JudgesCompetition{})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return domain.ErrAssignmentNotFound
		}
		return nil
	})
}

func (r *repository) ListJudges(page, limit int) ([]models.User, int64, error) {
	var total int64
	if err := r.db.Model(&models.User{}).
		Where("role = ?", string(types.UserRoleJudge)).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var judges []models.User
	err := r.db.
		Preload("Profile").
		Where("role = ?", string(types.UserRoleJudge)).
		Order("id").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&judges).Error
	return judges, total, err
}

func (r *repository) GetJudge(id uint) (*models.User, error) {
	var judge models.User
	err := r.db.Preload("Profile").
		Where("id = ? AND role = ?", id, string(types.UserRoleJudge)).
		First(&judge).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &judge, nil
}
