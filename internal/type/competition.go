package types

type CompetitionStatus string

const (
	CompetitionStatusDraft    CompetitionStatus = "draft"
	CompetitionStatusOpen     CompetitionStatus = "open"
	CompetitionStatusClosed   CompetitionStatus = "closed"
	CompetitionStatusFinished CompetitionStatus = "finished"
)

type PayoutStatus string

const (
	PayoutStatusPending   PayoutStatus = "pending"
	PayoutStatusCompleted PayoutStatus = "completed"
	PayoutStatusFailed    PayoutStatus = "failed"
)

type UserRole string

const (
	UserRoleParticipant UserRole = "participant"
	UserRoleJudge       UserRole = "judge"
	UserRoleOrganizer   UserRole = "organizer"
)
