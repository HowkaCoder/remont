package entity

import "gorm.io/gorm"

type Document struct {
	gorm.Model
	ID               uint   `gorm:"primaryKey" json:"id"`
	Name             string `gorm:"not null" json:"name"`
	Filepath         string `gorm:"not null" json:"filepath"`
	ProjectID        uint   `gorm:"not null" json:"projectID"`
	DocumentFolderID uint   `gorm:"not null" json:"document_folder_id`
}
