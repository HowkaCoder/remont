package usecase

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/repository"
)

type DocumentUsecase interface {
	GetAllDocuments() ([]entity.Document, error)
	GetDocumentByID(id uint) (*entity.Document, error)
	CreateDocument(doc *entity.Document) error
	UpdateDocument(doc *entity.Document, id uint) error
	//DeleteDocument(id uint) error
}

type documentUsecase struct {
	repo repository.DocumentRepository
}

func NewDocumentUsecase(repo repository.DocumentRepository) *documentUsecase {
	return &documentUsecase{repo: repo}
}

func (du *documentUsecase) GetAllDocuments() ([]entity.Document, error) {
	return du.repo.GetAllDocuments()
}

func (du *documentUsecase) GetDocumentByID(id uint) (*entity.Document, error) {
	return du.repo.GetDocumentByID(id)
}

func (du *documentUsecase) CreateDocument(doc *entity.Document) error {
	return du.repo.CreateDocument(doc)
}

func (du *documentUsecase) UpdateDocument(doc *entity.Document, id uint) error {
	return du.repo.UpdateDocument(doc, id)
}

//func (du *documentUsecase) DeleteDocument(id uint) error {
//	return du.repo.DeleteDocument(id)
//}
