package dto

import (
	types "godance/internal/type"
	"time"

	"github.com/shopspring/decimal"
)

type CreateCompetitionRequest struct {
	Title            string          `json:"title" binding:"required"`
	Event_date       time.Time       `json:"date" binding:"required"`
	ParticipantLimit uint            `json:"participant_limit" binding:"required,min=0"`
	EntryFee         decimal.Decimal `json:"entry_fee" binding:"required,min=0"`
}

type UpdateCompetitionRequest struct {
	Title            *string                  `json:"title,omitempty"`
	ParticipantLimit *int                     `json:"participant_limit,omitempty"`
	Status           *types.CompetitionStatus `json:"status,omitempty"`
}

type CompetitionItem struct {
	OrganizerID        uint     `json:"organizer_id"`
	Title              string   `json:"title"`
	EventDate          string   `json:"event_date"`
	ParticipantLimit   *int     `json:"participant_limit"`
	EntryFee           *float64 `json:"entry_fee"`
	FeedbackFeePercent *float64 `json:"organizer_fee_percent"`
	Status             string   `json:"status"`
}

type CompetitionListResponse struct {
	Data       []CompetitionItem `json:"data"`
	Pagination Pagination        `json:"pagination"`
}
