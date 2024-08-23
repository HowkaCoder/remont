package entity

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID      uint          `gorm:"primaryKey"`
	Title   string        `gorm:"not null" json:"title"`
	Members []ProjectRole `gorm:"foreignKey:ProjectID" json:"members"`
	Chars   []Char        `gorm:"foreignKey:ProjectID" json:"characteristics"`
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
