package entity

import (
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

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

	Roles []Role `gorm:"many2many:user_roles"`

	Projects []ProjectRole `gorm:"foreignKey:UserID"`
}

type Role struct {
	gorm.Model
	ID          uint         `gorm:"primaryKey"`
	Name        string       `json:"name"`
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

type JwtClaims struct {
	UserID    uint   `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Projects  []uint `json:"project"`
	jwt.StandardClaims
}
