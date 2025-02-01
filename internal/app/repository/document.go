package repository

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/gorm"
)

type DocumentRepository interface {
	GetAllDocuments() ([]entity.Document, error)
	GetDocumentByID(id uint) (*entity.Document, error)
	CreateDocument(doc *entity.Document) error
	UpdateDocument(doc *entity.Document, id uint) error
	GetDocumentsByFolderID(folderID uint) ([]entity.Document, error)
	DeleteDocument(id uint) error
}

type documentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(dbs *gorm.DB) *documentRepository {
	return &documentRepository{db: dbs}
}

func (dr *documentRepository) GetAllDocuments() ([]entity.Document, error) {
	var docs []entity.Document
	if err := dr.db.Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (dr *documentRepository) GetDocumentByID(id uint) (*entity.Document, error) {
	var doc *entity.Document
	if err := dr.db.First(&doc, id).Error; err != nil {
		return nil, err
	}
	return doc, nil
}

func (dr *documentRepository) CreateDocument(doc *entity.Document) error {
	return dr.db.Create(&doc).Error
}

func (dr *documentRepository) UpdateDocument(doc *entity.Document, id uint) error {
	var eDoc *entity.Document
	if err := dr.db.First(&eDoc, id).Error; err != nil {
		return err
	}

	if doc.Name != "" {
		eDoc.Name = doc.Name
	}
	if doc.Filepath != "" {
		eDoc.Filepath = doc.Filepath
	}
	if doc.ProjectID != 0 {
		eDoc.ProjectID = doc.ProjectID
	}
	if doc.DocumentFolderID != 0 {
		eDoc.DocumentFolderID = doc.DocumentFolderID
	}
	if doc.FileExt != "" {
		eDoc.FileExt = doc.FileExt
	}

	return dr.db.Save(&eDoc).Error
}

func (dr *documentRepository) DeleteDocument(id uint) error {
	var doc *entity.Document
	if err := dr.db.First(&doc, id).Error; err != nil {
		return err
	}
	return dr.db.Delete(&doc).Error
}

func (dr *documentRepository) GetDocumentsByFolderID(folderID uint) ([]entity.Document, error) {
	var docs []entity.Document
	err := dr.db.Where("document_folder_id =?", folderID).Find(&docs).Error
	if err != nil {
		return nil, err
	}
	return docs, nil
}
