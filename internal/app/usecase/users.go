package usecase

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/repository"
)

type UserUsecase interface {
	GetAllUsers() ([]entity.User , error) 
	GetUserByID(id uint) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User , id uint) error
	DeleteUser(id uint) error 
}

type userUsecase struct {
	repo 	repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) *userUsecase{
	return &userUsecase{repo:repo}
}

func (ur *userUsecase) GetAllUsers() ([]entity.User , error) {
	return ur.repo.GetAllUsers()
}

func (ur *userUsecase) CreateUser(user *entity.User) error {
	return ur.repo.CreateUser(user)
}

func (ur *userUsecase) GetUserByID(id uint) (*entity.User , error) {
	return ur.repo.GetUserByID(id)
}

func (ur *userUsecase) UpdateUser(user *entity.User , id uint) error {
	return ur.repo.UpdateUser(user , id)
}

func (ur *userUsecase) DeleteUser(id uint) error {
	return ur.repo.DeleteUser(id)
}
