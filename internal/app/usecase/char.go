package usecase

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/repository"
)

type CharUsecase interface {
	GetAllCharsByProjectID(id uint) ([]entity.Char, error)
	CreateChar(char *entity.Char) error
	UpdateChar(char *entity.Char, id uint) error
	DeleteChar(id uint) error
}

type charUsecase struct {
	repo repository.CharRepository
}

func NewCharUsecase(repo repository.CharRepository) *charUsecase {
	return &charUsecase{repo: repo}
}

func (cu *charUsecase) GetAllCharsByProjectID(id uint) ([]entity.Char, error) {
	return cu.repo.GetAllCharsByProjectID(id)
}

func (cu *charUsecase) CreateChar(char *entity.Char) error {
	return cu.repo.CreateChar(char)
}

func (cu *charUsecase) UpdateChar(char *entity.Char, id uint) error {
	return cu.repo.UpdateChar(char, id)
}

func (cu *charUsecase) DeleteChar(id uint) error {
	return cu.repo.DeleteChar(id)
}
