package entity

import "gorm.io/gorm" 

type Act struct {
  gorm.Model
  ID        uint        `gorm:primaryKey json:"id"`
  Name      string      `gorm:"not null" json:"name"`
  UnitPrice uint        `gorm:"not null" json:"unit_price"`
  Quantity  uint        `gorm:"not null" json:"quantity"`
  ProjectID uint        `gorm:"not null" json:"projectID"`
}
