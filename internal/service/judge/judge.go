package judge

import (
	"godance/internal/domain"
	"godance/internal/dto"
	"godance/internal/models"
)

const MaxJudgesPerCompetition = 5

type Service struct {
	repo domain.JudgeRepository
}

func NewService(repo domain.JudgeRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListForCompetition(competitionID uint) ([]dto.JudgeItem, error) {
	assignments, err := s.repo.ListByCompetition(competitionID)
	if err != nil {
		return nil, err
	}
	items := make([]dto.JudgeItem, len(assignments))
	for i, jc := range assignments {
		items[i] = toJudgeItem(jc)
	}
	return items, nil
}

func (s *Service) Assign(competitionID, judgeID, callerID uint) (*dto.JudgeItem, error) {
	assignment, err := s.repo.Assign(competitionID, judgeID, callerID, MaxJudgesPerCompetition)
	if err != nil {
		return nil, err
	}
	item := toJudgeItem(*assignment)
	return &item, nil
}

func toJudgeItem(jc models.JudgesCompetition) dto.JudgeItem {
	return dto.JudgeItem{
		ID:            jc.ID,
		CompetitionID: jc.CompetitionID,
		JudgeID:       jc.JudgeID,
		FullName:      jc.Judge.Profile.FullName,
		Rating:        jc.Judge.Profile.Rating,
		AssignedAt:    jc.AssignedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
