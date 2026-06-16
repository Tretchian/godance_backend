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
	ErrAssignmentNotFound  = errors.New("judge assignment not found")
	ErrJudgeHasRequests    = errors.New("judge has feedback requests in this competition")

	ErrAlreadyRegistered       = errors.New("participant already registered")
	ErrRegistrationClosed      = errors.New("registration is closed")
	ErrParticipantLimitReached = errors.New("participant limit reached")

	ErrFeedbackRequestNotFound  = errors.New("feedback request not found")
	ErrFeedbackResponseNotFound = errors.New("feedback response not found")
	ErrJudgeNotAssigned         = errors.New("judge is not assigned to this competition")
	ErrInvalidStatusTransition  = errors.New("operation not allowed in current status")
	ErrNotRequestParticipant    = errors.New("not the request participant")
	ErrNotRequestJudge          = errors.New("not the assigned judge of the request")
	ErrResponseAlreadyExists    = errors.New("feedback response already submitted")
	ErrRatingAlreadyExists      = errors.New("feedback already rated")

	ErrVideoNotFound           = errors.New("video not found")
	ErrNotVideoOwner           = errors.New("not the video owner")
	ErrRequestNotAwaitingVideo = errors.New("feedback request is not awaiting a video")
	ErrVideoObjectMissing      = errors.New("uploaded object not found in storage")
	ErrVideoGone               = errors.New("video has been deleted")

	ErrInvalidToken         = errors.New("invalid or expired token")
	ErrRegistrationNotFound = errors.New("registration not found")
)

const MaxJudgesPerCompetition = 5
