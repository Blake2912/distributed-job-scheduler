package models

import (
	"time"

	"gorm.io/gorm"
)

type JobExecution struct {
	gorm.Model

	JobID       uint   `gorm:"index"`
	Status      string `gorm:"type:varchar(50)"`
	StartedAt   time.Time
	CompletedAt time.Time
	Comments    string
	RetryCount  uint
}
