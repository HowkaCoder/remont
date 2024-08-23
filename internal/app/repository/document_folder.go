package repository

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/gorm"
)

type DocumentFolderRepository interface {
	GetAllDocumentFolders() ([]entity.DocumentFolder, error)
	GetDocumentFolderByID(id uint) (*entity.DocumentFolder, error)
	GetDocumentFoldersByProjectID(id uint) ([]entity.DocumentFolder, error)
	CreateDocumentFolder(folder *entity.DocumentFolder) error
	UpdateDocumentFolder(folder *entity.DocumentFolder, id uint) error
	DeleteDocumentFolder(id uint) error
}

type documentFolderRepository struct {
	db *gorm.DB
}

func NewDocumentFolderRepository(db *gorm.DB) *documentFolderRepository {
	return &documentFolderRepository{db: db}
}

func (dr *documentFolderRepository) GetAllDocumentFolders() ([]entity.DocumentFolder, error) {
	var folders []entity.DocumentFolder
	if err := dr.db.Preload("Documents").Find(&folders).Error; err != nil {
		return nil, err
	}
	return folders, nil
}

func (dr *documentFolderRepository) GetDocumentFolderByID(id uint) (*entity.DocumentFolder, error) {
	var doc *entity.DocumentFolder
	if err := dr.db.Preload("Documents").First(&doc, id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

func (dr *documentFolderRepository) GetDocumentFoldersByProjectID(id uint) ([]entity.DocumentFolder, error) {
	var folders []entity.DocumentFolder
	if err := dr.db.Where("project_id =?", id).Find(&folders).Error; err != nil {
		return nil, err
	}

	return folders, nil
}

func (dr *documentFolderRepository) CreateDocumentFolder(folder *entity.DocumentFolder) error {
	return dr.db.Create(folder).Error
}

func (dr *documentFolderRepository) UpdateDocumentFolder(folder *entity.DocumentFolder, id uint) error {
	var eFolder *entity.DocumentFolder
	if err := dr.db.First(&eFolder, id).Error; err != nil {
		return err
	}

	if folder.Title == "" {
		eFolder.Title = folder.Title
	}
	if folder.ProjectID == 0 {
		eFolder.ProjectID = folder.ProjectID
	}

	return dr.db.Save(eFolder).Error
}

func (dr *documentFolderRepository) DeleteDocumentFolder(id uint) error {
	var folder *entity.DocumentFolder
	if err := dr.db.First(&folder, id).Error; err != nil {
		return err
	}
	return dr.db.Delete(&folder).Error
}
