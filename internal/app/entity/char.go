package entity

import "gorm.io/gorm"

type Char struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey"`
	ProjectID uint   `gorm:"not null" json:"projectID"`
	Title     string `gorm:"not null" json:"title"`
	Desc      string `gorm:"not null" json:"desc"`
}
