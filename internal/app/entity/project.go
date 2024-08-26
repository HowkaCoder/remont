package entity

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID              uint             `gorm:"primaryKey"`
	Title           string           `gorm:"not null" json:"title"`
	Members         []ProjectRole    `gorm:"foreignKey:ProjectID" json:"members"`
	Chars           []Char           `gorm:"foreignKey:ProjectID" json:"characteristics"`
	PhotoFolders    []PhotoFolder    `gorm:"foreignKey:ProjectID" json:"photo_folders"`
	DocumentFolders []DocumentFolder `gorm:"foreignKey:ProjectID" json:"document_folders"`
	States          []State          `gorm:"foreignKey:ProjectID" json:"states"`
}

type ProjectRole struct {
	gorm.Model
	ProjectID uint
	Project   Project `gorm:"foreignKey:ProjectID"`
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	RoleID    uint
	Role      Role `gorm:"foreignKey:RoleID"`
}
