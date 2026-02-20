package models

import (
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"gorm.io/gorm"
)

type JobExecution struct {
	ID        uint      `gorm:"primaryKey"`
	JobID     uint      `gorm:"index:idx_job_exec_latest,priority:1"`
	CreatedAt time.Time `gorm:"index:idx_job_exec_latest,priority:2,sort:desc"`

	Status      database_constants.JobExecutionStatus `gorm:"type:varchar(50)"`
	StartedAt   *time.Time
	ScheduledAt time.Time
	RetryAt     *time.Time
	CompletedAt *time.Time
	Comments    string
	RetryCount  uint
	LeaseExpiry *time.Time
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Job Jobs `gorm:"foreignKey:JobID"`
}
