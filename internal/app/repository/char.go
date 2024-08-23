package repository

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/gorm"
)

type CharRepository interface {
	GetAllCharsByProjectID(id uint) ([]entity.Char, error)
	CreateChar(char *entity.Char) error
	UpdateChar(char *entity.Char, id uint) error
	DeleteChar(id uint) error
}

type charRepository struct {
	db *gorm.DB
}

func NewCharRepository(db *gorm.DB) *charRepository {
	return &charRepository{db: db}
}

func (cr *charRepository) GetAllCharsByProjectID(id uint) ([]entity.Char, error) {
	var chars []entity.Char
	if err := cr.db.Where("project_id = ?", id).Find(&chars).Error; err != nil {
		return nil, err
	}
	return chars, nil
}

func (cr *charRepository) CreateChar(char *entity.Char) error {
	return cr.db.Create(char).Error
}

func (cr *charRepository) UpdateChar(char *entity.Char, id uint) error {
	var eChar *entity.Char
	if err := cr.db.First(&eChar, id).Error; err != nil {
		return err
	}

	if char.Title != "" {
		eChar.Title = char.Title
	}
	if char.ProjectID != 0 {
		eChar.ProjectID = char.ProjectID
	}
	if char.Desc != "" {
		eChar.Desc = char.Desc
	}
	return cr.db.Save(&eChar).Error
}

func (cr *charRepository) DeleteChar(id uint) error {
	var char *entity.Char
	if err := cr.db.First(&char, id).Error; err != nil {
		return err
	}
	return cr.db.Delete(&char).Error
}
