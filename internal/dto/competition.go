package dto

import (
	"time"

	types "godance/internal/type"

	"github.com/shopspring/decimal"
)

type CreateCompetitionRequest struct {
	Title            string          `json:"title" binding:"required"`
	EventDate        time.Time       `json:"date" binding:"required"`
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

type CompetitionSummary struct {
	CompetiitonID       uint               `json:"competition_id"`
	TotalRequests       uint               `json:"total_requests"`
	Completed           uint               `json:"completed"`
	Pending             uint               `json:"pending"`
	Refunded            uint               `json:"refunded"`
	TotalCapturedAmount decimal.Decimal    `json:"total_captured_amount"`
	OrganizerShare      decimal.Decimal    `json:"organizer_share"`
	PayoutStatus        types.PayoutStatus `json:"payout_status"`
}
