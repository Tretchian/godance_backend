package models

import (
	"time"

	"gorm.io/gorm"
)

type JudgesCompetition struct {
	gorm.Model
	CompetitionID uint        `gorm:"not null;uniqueIndex:idx_comp_judge;index"`
	Competition   Competition `gorm:"foreignKey:CompetitionID"`
	JudgeID       uint        `gorm:"not null;uniqueIndex:idx_comp_judge"`
	Judge         User        `gorm:"foreignKey:JudgeID"`
	AssignedAt    time.Time   `gorm:"autoCreateTime"`
}
