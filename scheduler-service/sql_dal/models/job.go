package models

import (
	"database/sql"

	"gorm.io/gorm"
)

type Jobs struct {
	gorm.Model

	Name                    string `gorm:"type:varchar(500)"`
	Type                    string `gorm:"type:varchar(500)"`
	Config                  string `gorm:"type:varchar(500)"`
	Enabled                 bool
	NextRunAt               sql.NullTime
	ShouldRetryAfterBackoff bool

	JobExecutions []JobExecution `gorm:"foreignKey:JobID;references:ID"`
}
