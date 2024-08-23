package usecase

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/repository"
)

type DocumentFolderUsecase interface {
	GetAllDocumentFolders() ([]entity.DocumentFolder, error)
	GetDocumentFolderByID(id uint) (*entity.DocumentFolder, error)
	GetDocumentFoldersByProjectID(id uint) ([]entity.DocumentFolder, error)
	CreateDocumentFolder(folder *entity.DocumentFolder) error
	UpdateDocumentFolder(folder *entity.DocumentFolder, id uint) error
	DeleteDocumentFolder(id uint) error
}

type documentFolderUsecase struct {
	repo repository.DocumentFolderRepository
}

func NewDocumentFolderUsecase(repo repository.DocumentFolderRepository) *documentFolderUsecase {
	return &documentFolderUsecase{repo: repo}
}

func (du *documentFolderUsecase) GetAllDocumentFolders() ([]entity.DocumentFolder, error) {
	return du.repo.GetAllDocumentFolders()
}

func (du *documentFolderUsecase) GetDocumentFolderByID(id uint) (*entity.DocumentFolder, error) {
	return du.repo.GetDocumentFolderByID(id)
}

func (du *documentFolderUsecase) GetDocumentFoldersByProjectID(id uint) ([]entity.DocumentFolder, error) {
	return du.repo.GetDocumentFoldersByProjectID(id)
}

func (du *documentFolderUsecase) CreateDocumentFolder(folder *entity.DocumentFolder) error {
	return du.repo.CreateDocumentFolder(folder)
}

func (du *documentFolderUsecase) UpdateDocumentFolder(folder *entity.DocumentFolder, id uint) error {
	return du.repo.UpdateDocumentFolder(folder, id)
}

func (du *documentFolderUsecase) DeleteDocumentFolder(id uint) error {
	return du.repo.DeleteDocumentFolder(id)
}
