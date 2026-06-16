package dto

type AssignJudgeRequest struct {
	JudgeID uint `json:"judge_id" binding:"required"`
}

type JudgeItem struct {
	ID            uint    `json:"id"`
	CompetitionID uint    `json:"competition_id"`
	JudgeID       uint    `json:"judge_id"`
	FullName      string  `json:"full_name,omitempty"`
	Rating        float64 `json:"rating,omitempty"`
	AssignedAt    string  `json:"assigned_at"`
}

type JudgeProfile struct {
	ID       uint    `json:"id"`
	FullName string  `json:"full_name"`
	Bio      string  `json:"bio"`
	Rating   float64 `json:"rating"`
}

type JudgeProfileList struct {
	Data       []JudgeProfile `json:"data"`
	Pagination Pagination     `json:"pagination"`
}
