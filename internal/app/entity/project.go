package entity

import (
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	ID      uint `gorm:"primaryKey"`
	Title   string
	Members []ProjectRole
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
