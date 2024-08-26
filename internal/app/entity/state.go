package entity

import (
	"gorm.io/gorm"
)

type State struct {
	gorm.Model
	ID             uint   `gorm:"primaryKey"`
	Position       uint   `gorm:"not null" json:"position"`
	ProjectID      uint   `gorm:"not null" json:"project_id"`
	Job            string `gorm:"not null" json:"job"`
	NecessaryMoney uint   `gorm:"not null" json:"necessary_money"`
	PaidMoney      uint   `gorm:"not null" json:"paid_money"`
	Deadline       string `gorm:"not null" json:"deadline"`
	Status         string `gorm:"not null" json:"status"`
	DonePart       uint   `gorm:"not null" json:"done_part"`
	Workers        []User `gorm:"many2many:state_users" json:"workers"`
}

type StateUser struct {
	gorm.Model
	ID      uint `gorm:"primaryKey"`
	StateID uint `gorm:"not null" json:"state_id"`
	UserID  uint `gorm:"not null" json:"user_id"`
}
