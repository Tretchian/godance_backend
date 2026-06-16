package domain

import "godance/internal/models"

type JudgeRepository interface {
	ListByCompetition(competitionID uint) ([]models.JudgesCompetition, error)
	Assign(competitionID, judgeID, callerID uint, limit int) (*models.JudgesCompetition, error)
	// Unassign снимает судью с соревнования. Запрещено (ErrJudgeHasRequests),
	// если по судье уже есть запросы ОС в этом соревновании.
	Unassign(competitionID, judgeID, callerID uint) error
	// ListJudges возвращает каталог пользователей-судей с пагинацией.
	ListJudges(page, limit int) ([]models.User, int64, error)
	GetJudge(id uint) (*models.User, error)
}
