package domain

import "errors"

var (
	ErrCompetitionNotFound = errors.New("competition not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrNotOwner            = errors.New("not the competition organizer")
	ErrNotJudge            = errors.New("user is not a judge")
	ErrJudgeLimitReached   = errors.New("judge limit reached")
	ErrAlreadyAssigned     = errors.New("judge already assigned")
	ErrCompetitionClosed   = errors.New("competition not open for changes")

	ErrAlreadyRegistered       = errors.New("participant already registered")
	ErrRegistrationClosed      = errors.New("registration is closed")
	ErrParticipantLimitReached = errors.New("participant limit reached")
)

const MaxJudgesPerCompetition = 5
