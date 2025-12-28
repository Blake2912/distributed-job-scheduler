package datamodels

import (
	"time"

	"gorm.io/gorm"
)

type JobExecution struct {
	gorm.Model
	ID          uint
	JobID       uint
	Status      string `gorm:"type:varchar(50)"`
	StartedAt   time.Time
	CompletedAt time.Time
	Comments    string
	RetryCount  uint
}
