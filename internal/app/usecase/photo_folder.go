package usecase

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/repository"
)

type PhotoFolderUseCase interface {
	GetAllPhotoFolders() ([]entity.PhotoFolder, error)
	GetPhotoFolderByID(id uint) (*entity.PhotoFolder, error)
	GetPhotoFoldersByProjectID(userID uint) ([]entity.PhotoFolder, error)
	CreatePhotoFolder(folder *entity.PhotoFolder) error
	UpdatePhotoFolder(folder *entity.PhotoFolder, id uint) error
	DeletePhotoFolder(id uint) error
}

type photoFolderUsecase struct {
	repo repository.PhotoFolder
}

func NewPhotoFolderUsecase(repo repository.PhotoFolder) *photoFolderUsecase {
	return &photoFolderUsecase{repo: repo}
}

func (pu *photoFolderUsecase) GetAllPhotoFolders() ([]entity.PhotoFolder, error) {
	return pu.repo.GetAllPhotoFolders()
}

func (pu *photoFolderUsecase) GetPhotoFolderByID(id uint) (*entity.PhotoFolder, error) {
	return pu.repo.GetPhotoFolderByID(id)
}

func (pu *photoFolderUsecase) GetPhotoFoldersByProjectID(userID uint) ([]entity.PhotoFolder, error) {
	return pu.repo.GetPhotoFoldersByProjectID(userID)
}

func (pu *photoFolderUsecase) CreatePhotoFolder(folder *entity.PhotoFolder) error {
	return pu.repo.CreatePhotoFolder(folder)
}

func (pu *photoFolderUsecase) UpdatePhotoFolder(folder *entity.PhotoFolder, id uint) error {
	return pu.repo.UpdatePhotoFolder(folder, id)
}

func (pu *photoFolderUsecase) DeletePhotoFolder(id uint) error {
	return pu.repo.DeletePhotoFolder(id)
}
