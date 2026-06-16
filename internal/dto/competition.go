package dto

import (
	types "godance/internal/type"

	"github.com/shopspring/decimal"
)

type CreateCompetitionRequest struct {
	Title            string          `json:"title" binding:"required"`
	EventDate        Date            `json:"event_date" binding:"required"`
	ParticipantLimit *int            `json:"participant_limit" binding:"omitempty,min=1"`
	EntryFee         decimal.Decimal `json:"entry_fee" binding:"required"`
}

type UpdateCompetitionRequest struct {
	Title            *string                  `json:"title,omitempty"`
	ParticipantLimit *int                     `json:"participant_limit,omitempty"`
	Status           *types.CompetitionStatus `json:"status,omitempty" binding:"omitempty,oneof=draft open closed finished"`
}

// Competition — полная схема соревнования (ответ GET/POST/PATCH).
type Competition struct {
	ID               uint     `json:"id"`
	OrganizerID      uint     `json:"organizer_id"`
	OrganizerName    string   `json:"organizer_name"`
	OrganizerSurname string   `json:"organizer_surname"`
	Title            string   `json:"title"`
	EventDate        string   `json:"event_date"`
	ParticipantLimit *int     `json:"participant_limit"`
	EntryFee         *float64 `json:"entry_fee"`
	Status           string   `json:"status"`
	PayoutStatus     string   `json:"payout_status"`
	CreatedAt        string   `json:"created_at"`
}

type CompetitionListResponse struct {
	Data       []Competition `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

type CompetitionSummary struct {
	CompetitionID       uint    `json:"competition_id"`
	TotalRequests       int64   `json:"total_requests"`
	Completed           int64   `json:"completed"`
	Pending             int64   `json:"pending"`
	Refunded            int64   `json:"refunded"`
	TotalCapturedAmount float64 `json:"total_captured_amount"`
	OrganizerShare      float64 `json:"organizer_share"`
	PayoutStatus        string  `json:"payout_status"`
}
