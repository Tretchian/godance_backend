package models

import (
	"time"

	"gorm.io/gorm"
)

type FeedbackRating struct {
	gorm.Model
	ResponseID    uint      `gorm:"uniqueIndex;not null"`
	ParticipantID uint      `gorm:"not null"`
	Participant   User      `gorm:"foreignKey:ParticipantID"`
	Score         int16     `gorm:"type:smallint;not null"` // 1–5
	RatedAt       time.Time `gorm:"autoCreateTime"`
}

type FeedbackRequest struct {
	gorm.Model
	JudgeID            uint `gorm:"not null;index:idx_feedback_requests_judge_status"`
	Judge              User `gorm:"foreignKey:JudgeID"`
	ParticipantID      uint `gorm:"not null"`
	Participant        User `gorm:"foreignKey:ParticipantID"`
	VideoID            uint `gorm:"not null"`
	Video              Video
	Price              *float64   `gorm:"type:numeric(10,2)"`
	ParticipantComment *string    `gorm:"type:text"`
	Status             string     `gorm:"type:varchar(20);index:idx_feedback_requests_judge_status"` // pending | completed | refunded
	DeadlineAt         *time.Time `gorm:"type:timestamptz"`                                          // created_at + 30 days

	Response *FeedbackResponse `gorm:"foreignKey:RequestID"`
}

type FeedbackResponse struct {
	gorm.Model
	RequestID       uint   `gorm:"uniqueIndex;not null"`
	Strengths       string `gorm:"type:text;not null"`
	Errors          string `gorm:"type:text;not null"`
	Recommendations string `gorm:"type:text;not null"`

	Rating *FeedbackRating `gorm:"foreignKey:ResponseID"`
}
