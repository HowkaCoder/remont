package repository

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/gorm"
)

type PhotoFolder interface {
	GetAllPhotoFolders() ([]entity.PhotoFolder, error)
	GetPhotoFolderByID(id uint) (*entity.PhotoFolder, error)
	GetPhotoFoldersByProjectID(id uint) ([]entity.PhotoFolder, error)
	CreatePhotoFolder(folder *entity.PhotoFolder) error
	UpdatePhotoFolder(folder *entity.PhotoFolder, id uint) error
	DeletePhotoFolder(id uint) error
}

type PhotoFolderRepository struct {
	db *gorm.DB
}

func NewPhotoFolderRepository(db *gorm.DB) *PhotoFolderRepository {
	return &PhotoFolderRepository{db: db}
}

func (pr *PhotoFolderRepository) GetAllPhotoFolders() ([]entity.PhotoFolder, error) {
	var folders []entity.PhotoFolder
	if err := pr.db.Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

func (pr *PhotoFolderRepository) GetPhotoFolderByID(id uint) (*entity.PhotoFolder, error) {
	var folder *entity.PhotoFolder
	if err := pr.db.First(&folder, id).Error; err != nil {
		return nil, err
	}
	return folder, nil
}

func (pr *PhotoFolderRepository) GetPhotoFoldersByProjectID(id uint) ([]entity.PhotoFolder, error) {
	var folders []entity.PhotoFolder
	if err := pr.db.Where("project_id =?", id).Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

func (pr *PhotoFolderRepository) CreatePhotoFolder(folder *entity.PhotoFolder) error {
	return pr.db.Create(folder).Error
}

func (pr *PhotoFolderRepository) UpdatePhotoFolder(folder *entity.PhotoFolder, id uint) error {
	var eFolder *entity.PhotoFolder
	if err := pr.db.First(&eFolder, id).Error; err != nil {
		return err
	}
	return pr.db.Save(folder).Error
}

func (pr *PhotoFolderRepository) DeletePhotoFolder(id uint) error {
	var folder *entity.PhotoFolder
	if err := pr.db.First(&folder, id).Error; err != nil {
		return err
	}
	return pr.db.Delete(&folder).Error
}
