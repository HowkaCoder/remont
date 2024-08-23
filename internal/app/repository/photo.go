package repository

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/gorm"
)

type PhotoRepository interface {
	GetAllPhotos() ([]entity.Photo, error)
	GetPhotoByID(id uint) (*entity.Photo, error)
	CreatePhoto(photo *entity.Photo) error
	GetPhotosByFolderID(folderID uint) ([]entity.Photo, error)
	UpdatePhoto(photo *entity.Photo, id uint) error
	DeletePhoto(id uint) error
}

type photoRepository struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) *photoRepository {
	return &photoRepository{db: db}
}

func (pr *photoRepository) GetAllPhotos() ([]entity.Photo, error) {
	var photos []entity.Photo
	if err := pr.db.Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}

func (pr *photoRepository) GetPhotoByID(id uint) (*entity.Photo, error) {
	var photo *entity.Photo
	if err := pr.db.First(&photo, id).Error; err != nil {
		return nil, err
	}
	return photo, nil
}

func (pr *photoRepository) CreatePhoto(photo *entity.Photo) error {
	return pr.db.Create(photo).Error
}

func (pr *photoRepository) UpdatePhoto(photo *entity.Photo, id uint) error {
	var ePhoto *entity.Photo
	if err := pr.db.First(&ePhoto, id).Error; err != nil {
		return err
	}

	if photo.Title != "" {
		ePhoto.Title = photo.Title
	}
	if photo.ProjectID != 0 {
		ePhoto.ProjectID = photo.ProjectID
	}
	if photo.Filepath != "" {
		ePhoto.Filepath = photo.Filepath
	}
	if photo.PhotoFolderID != 0 {
		ePhoto.PhotoFolderID = photo.PhotoFolderID
	}

	return pr.db.Save(ePhoto).Error
}

func (pr *photoRepository) DeletePhoto(id uint) error {
	var photo *entity.Photo
	if err := pr.db.First(&photo, id).Error; err != nil {
		return err
	}
	return pr.db.Delete(&photo).Error
}

func (pr *photoRepository) GetPhotosByFolderID(folderID uint) ([]entity.Photo, error) {
	var photos []entity.Photo
	if err := pr.db.Where("photo_folder_id =?", folderID).Find(&photos).Error; err != nil {
		return nil, err
	}
	return photos, nil
}
