package dto

// ── Requests (входящие) ───────────────────────────────────────────────────

type CreateFeedbackRequest struct {
	JudgeID       uint   `json:"judge_id" binding:"required"`
	CompetitionID uint   `json:"competition_id" binding:"required"`
	Comment       string `json:"comment" binding:"required"`
}

type FeedbackRequestAction struct {
	Action string `json:"action" binding:"required,oneof=confirm"`
}

type CreateFeedbackResponse struct {
	RequestID       uint   `json:"request_id" binding:"required"`
	Strengths       string `json:"strengths" binding:"required,min=10"`
	Errors          string `json:"errors" binding:"required,min=10"`
	Recommendations string `json:"recommendations" binding:"required,min=10"`
}

type CreateFeedbackRating struct {
	ResponseID uint `json:"response_id" binding:"required"`
	Score      int8 `json:"score" binding:"required,min=1,max=5"`
}

type PaymentWebhook struct {
	Status     string   `json:"status" binding:"required"`
	RelatedID  uint     `json:"related_id" binding:"required"`
	HoldID     *string  `json:"hold_id"`
	Amount     *float64 `json:"amount"`
	GatewayRef string   `json:"gateway_ref"`
}

// ── Responses (исходящие) ─────────────────────────────────────────────────

type FeedbackRequestCreated struct {
	RequestID  uint   `json:"request_id"`
	PaymentURL string `json:"payment_url"`
	Status     string `json:"status"`
}

type FeedbackRequestItem struct {
	ID            uint     `json:"id"`
	JudgeID       uint     `json:"judge_id"`
	ParticipantID uint     `json:"participant_id"`
	VideoID       *uint    `json:"video_id"`
	Comment       *string  `json:"comment"`
	Price         *float64 `json:"price"`
	Status        string   `json:"status"`
	CreatedAt     string   `json:"created_at"`
	DeadlineAt    *string  `json:"deadline_at"`
}

type FeedbackResponseItem struct {
	ID              uint   `json:"id"`
	RequestID       uint   `json:"request_id"`
	Strengths       string `json:"strengths"`
	Errors          string `json:"errors"`
	Recommendations string `json:"recommendations"`
	SubmittedAt     string `json:"submitted_at"`
}

type FeedbackRatingItem struct {
	ID             uint    `json:"id"`
	ResponseID     uint    `json:"response_id"`
	Score          int8    `json:"score"`
	JudgeNewRating float64 `json:"judge_new_rating"`
	RatedAt        string  `json:"rated_at"`
}
