package datamodels

import (
	"time"

	databaseconstants "example.com/constants/database_constants"
	"gorm.io/gorm"
)

// TODO:: This is a bad approach on maintaining the images in longer run,
// we need to implement a CI Pipeline which will publish the images and we make use of the stable image
// improve this later.
type ImageInformation struct {
	gorm.Model
	ID           uint
	Image        string
	ImageVersion string
	Type         databaseconstants.ImageType `gorm:"type:varchar(50);not null"`
	CreatedAt    time.Time
}
