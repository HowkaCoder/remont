package entity

import "gorm.io/gorm"

type Photo struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	Title         string `gorm:"null" json:"title"`
	PhotoFolderID uint   `gorm:"not null" json:"folderID"`
	ProjectID     uint   `gorm:"not null" json:"projectID"`
	Filepath      string `gorm:"null" json:"filePath"`
}
