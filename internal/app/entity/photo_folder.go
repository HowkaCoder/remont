package entity

import (
	"gorm.io/gorm"
)

type PhotoFolder struct {
	gorm.Model
	ID        uint    `gorm:"primaryKey"`
	Title     string  `gorm:"null" json:"title"`
	ProjectID uint    `gorm:"not null" json:"projectID"`
	Photos    []Photo `gorm:"foreignKey:PhotoFolderID"`
}
