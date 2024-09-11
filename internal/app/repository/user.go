package repository

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	GetAllUsers()	([]entity.User , error)
	GetUserByID(id uint) (*entity.User , error)
	CreateUser(user *entity.User) error 
	UpdateUser(user *entity.User , id uint) error 
	DeleteUser(id uint) error 



}

type userRepository struct {
	db 	*gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository{
	return &userRepository{db:db}
}

func (ur *userRepository) GetAllUsers() ([]entity.User , error) {
	var users []entity.User
	if err := ur.db.Find(&users).Error; err != nil {
		return nil,err
	}

	return users,nil
}

func (ur *userRepository) GetUserByID(id uint) (*entity.User , error) {
	var user *entity.User
	if err := ur.db.First(&user , id).Error; err != nil {
		return nil,err
	}

	return user,nil
}

func (ur *userRepository) CreateUser(user *entity.User) error {
	if err := ur.db.Create(user).Error; err != nil {
		return err
	}
	return nil
} 

func (ur *userRepository) UpdateUser(user *entity.User , id uint) error {
	var eUser *entity.User
	if err := ur.db.First(&eUser , id).Error; err != nil {
		return err
	}

	if user.FirstName != "" {
		eUser.FirstName = user.FirstName
	} else if user.MiddleName != "" {
		eUser.MiddleName = user.MiddleName
	} else if user.LastName != "" {
		eUser.LastName = user.LastName
	} else if user.Email != "" {
		eUser.Email = user.Email 
	} else if user.PhoneNumber != "" {
		eUser.PhoneNumber = user.PhoneNumber
	} else if user.Password != "" {
		eUser.Password = user.Password
	}

	if err := ur.db.Save(&eUser).Error; err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) DeleteUser(id uint) error {
	var user *entity.User
	if err := ur.db.First(&user , id).Error; err != nil {
		return err
	}

	if err := ur.db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}


