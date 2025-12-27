package datamodels

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Jobs struct {
	gorm.Model
	ID uint  
	Name string
	Type string
	Config string
	Enabled bool
	NextRunAt sql.NullTime
	CreatedAt time.Time
	ShouldRetryAfterBackoff bool
}