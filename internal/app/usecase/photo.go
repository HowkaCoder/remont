package usecase

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/repository"
)

type PhotoUsecase interface {
	GetAllPhotos() ([]entity.Photo, error)
	GetPhotoByID(id uint) (*entity.Photo, error)
	CreatePhoto(photo *entity.Photo) error
	GetPhotosByFolderID(folderID uint) ([]entity.Photo, error)
	UpdatePhoto(photo *entity.Photo, id uint) error
	DeletePhoto(id uint) error
}

type photoUsecase struct {
	repo repository.PhotoRepository
}

func NewPhotoUsecase(repo repository.PhotoRepository) *photoUsecase {
	return &photoUsecase{repo: repo}
}

func (u *photoUsecase) GetAllPhotos() ([]entity.Photo, error) {
	return u.repo.GetAllPhotos()
}

func (u *photoUsecase) GetPhotoByID(id uint) (*entity.Photo, error) {
	return u.repo.GetPhotoByID(id)

}

func (u *photoUsecase) CreatePhoto(photo *entity.Photo) error {
	return u.repo.CreatePhoto(photo)
}

func (u *photoUsecase) UpdatePhoto(photo *entity.Photo, id uint) error {
	return u.repo.UpdatePhoto(photo, id)
}

func (u *photoUsecase) DeletePhoto(id uint) error {
	return u.repo.DeletePhoto(id)
}

func (u *photoUsecase) GetPhotosByFolderID(folderID uint) ([]entity.Photo, error) {
	return u.repo.GetPhotosByFolderID(folderID)
}
