package entity

import "gorm.io/gorm"

type DocumentFolder struct {
	gorm.Model
	ID        uint       `gorm:"primaryKey"`
	Title     string     `gorm:"not null" json:"title"`
	ProjectID uint       `gorm:"not null" json:"projectID"`
	Documents []Document `gorm:"foreignKey:DocumentFolderID"`
}
