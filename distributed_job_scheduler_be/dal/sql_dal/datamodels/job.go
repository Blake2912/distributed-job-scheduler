package datamodels

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Jobs struct {
	gorm.Model
	ID                      uint
	Name                    string `gorm:"type:varchar(500)"`
	Type                    string `gorm:"type:varchar(500)"`
	Config                  string `gorm:"type:varchar(500)"`
	Enabled                 bool
	NextRunAt               sql.NullTime
	CreatedAt               time.Time
	ShouldRetryAfterBackoff bool
	JobExecutions []JobExecution
}
