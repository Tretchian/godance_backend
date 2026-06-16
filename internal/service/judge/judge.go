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

func (s *Service) Unassign(competitionID, judgeID, callerID uint) error {
	return s.repo.Unassign(competitionID, judgeID, callerID)
}

func (s *Service) ListCatalog(page, limit int) (*dto.JudgeProfileList, error) {
	judges, total, err := s.repo.ListJudges(page, limit)
	if err != nil {
		return nil, err
	}
	items := make([]dto.JudgeProfile, len(judges))
	for i, j := range judges {
		items[i] = toJudgeProfile(j)
	}
	return &dto.JudgeProfileList{
		Data: items,
		Pagination: dto.Pagination{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}, nil
}

func (s *Service) GetProfile(id uint) (*dto.JudgeProfile, error) {
	judge, err := s.repo.GetJudge(id)
	if err != nil {
		return nil, err
	}
	profile := toJudgeProfile(*judge)
	return &profile, nil
}

func toJudgeProfile(u models.User) dto.JudgeProfile {
	profile := dto.JudgeProfile{ID: u.ID}
	if u.Profile != nil {
		profile.FullName = u.Profile.FullName
		profile.Bio = u.Profile.Bio
		profile.Rating = u.Profile.Rating
	}
	return profile
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
