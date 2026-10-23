package repository

import (
	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/gorm"
)

const (
	ProjectRoleWorker  = 1
	ProjectRoleManager = 2
)

type AuthRepository interface {
	FindUserByEmail(email string) (*entity.User, error)
	FindUserWithRoles(id uint) (*entity.User, error)
	CreateUserWithRole(user *entity.User, roleName string) error
	ProjectIDs(userID uint) ([]uint, error)
	ProjectMembers(projectID uint) ([]entity.User, *entity.User, error)

	CreateRole(role *entity.Role) error
	GetRoles() ([]entity.Role, error)
	CreatePermission(permission *entity.Permission) error
	GetPermissions() ([]entity.Permission, error)
	AssignPermissionToRole(roleID, permissionID uint) error
	AssignRoleToUser(userID, roleID uint) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Preload("Roles").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FindUserWithRoles(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.Preload("Roles.Permissions").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) CreateUserWithRole(user *entity.User, roleName string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var role entity.Role
		if err := tx.Where("name = ?", roleName).First(&role).Error; err != nil {
			return err
		}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		return tx.Create(&entity.UserRole{UserID: user.ID, RoleID: role.ID}).Error
	})
}

func (r *authRepository) ProjectIDs(userID uint) ([]uint, error) {
	var ids []uint
	err := r.db.Model(&entity.ProjectRole{}).Where("user_id = ?", userID).Pluck("project_id", &ids).Error
	return ids, err
}

func (r *authRepository) ProjectMembers(projectID uint) ([]entity.User, *entity.User, error) {
	var roles []entity.ProjectRole
	if err := r.db.Where("project_id = ?", projectID).Find(&roles).Error; err != nil {
		return nil, nil, err
	}

	var workerIDs []uint
	var managerID uint
	for _, pr := range roles {
		switch pr.RoleID {
		case ProjectRoleWorker:
			workerIDs = append(workerIDs, pr.UserID)
		case ProjectRoleManager:
			managerID = pr.UserID
		}
	}

	workers := []entity.User{}
	if len(workerIDs) > 0 {
		if err := r.db.Where("id IN ?", workerIDs).Find(&workers).Error; err != nil {
			return nil, nil, err
		}
	}

	if managerID == 0 {
		return workers, nil, nil
	}
	var manager entity.User
	if err := r.db.First(&manager, managerID).Error; err != nil {
		return nil, nil, err
	}
	return workers, &manager, nil
}

func (r *authRepository) CreateRole(role *entity.Role) error {
	return r.db.Create(role).Error
}

func (r *authRepository) GetRoles() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *authRepository) CreatePermission(permission *entity.Permission) error {
	return r.db.Create(permission).Error
}

func (r *authRepository) GetPermissions() ([]entity.Permission, error) {
	var permissions []entity.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

func (r *authRepository) AssignPermissionToRole(roleID, permissionID uint) error {
	if err := r.db.First(&entity.Role{}, roleID).Error; err != nil {
		return err
	}
	if err := r.db.First(&entity.Permission{}, permissionID).Error; err != nil {
		return err
	}
	link := entity.RolePermission{RoleID: roleID, PermissionID: permissionID}
	return r.db.Where(link).FirstOrCreate(&link).Error
}

func (r *authRepository) AssignRoleToUser(userID, roleID uint) error {
	if err := r.db.First(&entity.User{}, userID).Error; err != nil {
		return err
	}
	if err := r.db.First(&entity.Role{}, roleID).Error; err != nil {
		return err
	}
	link := entity.UserRole{UserID: userID, RoleID: roleID}
	return r.db.Where(link).FirstOrCreate(&link).Error
}
