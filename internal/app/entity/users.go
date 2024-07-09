package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID            uint `gorm:"primaryKey"`
	FirstName     string
	LastName      string
	MiddleName    string
	Email         string `gorm:"unique"`
	Password      string
	PhoneNumber   string
	CreatedAt     string
	UserType      string
	DeactivatedAt *string
	Roles         []Role `gorm:"many2many:user_roles"`
}

type Role struct {
	gorm.Model
	ID          uint         `gorm:"primaryKey"`
	Name        string       `gorm:"unique"`
	Permissions []Permission `gorm:"many2many:role_permissions"`
}

type Permission struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type UserRole struct {
	gorm.Model
	UserID uint
	RoleID uint
}

type RolePermission struct {
	gorm.Model
	RoleID       uint
	PermissionID uint
}

type ProjectManager struct {
	ProjectID uint
	ManagerID uint
}

type ProjectWorker struct {
	ProjectID uint
	WorkerID  uint
}
