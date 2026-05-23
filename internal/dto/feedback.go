package dto

type CreateFeedbackRating struct {
	ResponseID uint `json:"response_id" binding:"required"`
	Score      int8 `json:"score" binding:"required,min=0,max=10"`
}

type CreateFeedbackRequest struct {
	JudgeID       uint   `json:"judge_id" binding:"required"`
	CompetitionID uint   `json:"competition_id" binding:"required"`
	Comment       string `json:"comment,omitempty" binding:"omitempty"`
}

type CreateFeedbackResponse struct {
	RequestID       uint   `json:"request_id" binding:"required"`
	Strengths       string `json:"strengths" binding:"required,min=10"`
	Errors          string `json:"errors" binding:"required,min=10"`
	Recommendations string `json:"recommendations" binding:"required,min=10"`
}
