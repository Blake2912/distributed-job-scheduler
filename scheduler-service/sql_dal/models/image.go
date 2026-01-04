package models

import (
	databaseconstants "github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"gorm.io/gorm"
)

// TODO:: This is a bad approach on maintaining the images in longer run,
// we need to implement a CI Pipeline which will publish the images and we make use of the stable image
// improve this later.
type ImageInformation struct {
	gorm.Model
	Image        string
	ImageVersion string                      `gorm:"index:idx_type_version"`
	ImageType    databaseconstants.ImageType `gorm:"type:varchar(50);not null;index:idx_type_version"`
}
