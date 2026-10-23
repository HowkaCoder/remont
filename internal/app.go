package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	PermManageUsers = "manage_users"
	PermManageRoles = "manage_roles"
)

var rolePermissions = map[string][]string{
	"worker": {
		"get_documents", "create_document", "update_document", "get_document_by_id", "get_documents_by_folder_id",
		"delete_document", "get_all_projects", "create_project", "update_project", "get_project_by_id", "delete_project",
		"get_all_project_roles", "create_project_role", "update_project_role", "get_project_role_by_id",
		"delete_project_role", "get_all_photos", "create_photo", "get_photos_by_folder_id", "get_photo_by_id",
		"delete_photo", "update_photo", "get_all_document_folders", "create_document_folder",
		"update_document_folder", "delete_document_folder", "get_photo_folder_by_project_id", "create_photo_folder",
		"update_photo_folder", "delete_photo_folder", "get_chars_by_project_id", "create_char", "update_char", "delete_char",
		"get_states_by_project_id", "create_state", "update_state", "delete_state", "create_state_relation", "delete_state_relation",
		"get_projects_by_worker_id", "get_projects_by_client_id", "create_detail", "update_detail", "get_detail_by_id", "delete_detail",
		"get_states_by_worker_id",
	},
	"manager": {
		"get_documents", "create_document", "update_document", "get_document_by_id", "get_documents_by_folder_id",
		"delete_document", "get_all_projects", "create_project", "update_project", "get_project_by_id", "delete_project",
		"get_all_project_roles", "create_project_role", "update_project_role", "get_project_role_by_id",
		"delete_project_role", "get_all_photos", "create_photo", "get_photos_by_folder_id", "get_photo_by_id",
		"delete_photo", "update_photo", "get_all_document_folders", "create_document_folder",
		"update_document_folder", "delete_document_folder", "get_photo_folder_by_project_id", "create_photo_folder",
		"update_photo_folder", "delete_photo_folder", "get_chars_by_project_id", "create_char", "update_char", "delete_char",
		"get_states_by_project_id", "create_state", "update_state", "delete_state", "create_state_relation", "delete_state_relation",
		"get_projects_by_worker_id", "get_projects_by_client_id", "create_detail", "update_detail", "get_detail_by_id", "delete_detail",
	},
	"client": {
		"get_documents", "get_document_by_id", "get_documents_by_folder_id", "get_all_projects", "get_project_by_id",
		"get_all_project_roles", "get_project_role_by_id", "get_all_photos", "get_photos_by_folder_id", "get_photo_by_id",
		"get_all_document_folders", "get_document_folder_by_folder_id", "get_document_folder_by_project_id",
		"get_all_photo_folders", "get_photo_folder_by_folder_id", "get_photo_folder_by_project_id",
		"get_chars_by_project_id", "get_states_by_project_id", "get_projects_by_worker_id", "get_projects_by_client_id", "create_detail", "update_detail", "get_detail_by_id", "delete_detail",
	},
	"prorab": {
		"get_documents", "create_document", "update_document", "get_document_by_id", "get_documents_by_folder_id",
		"delete_document", "get_all_projects", "create_project", "update_project", "get_project_by_id", "delete_project",
		"get_all_project_roles", "create_project_role", "update_project_role", "get_project_role_by_id",
		"delete_project_role", "get_all_photos", "create_photo", "get_photos_by_folder_id", "get_photo_by_id",
		"delete_photo", "update_photo", "get_all_document_folders", "create_document_folder",
		"update_document_folder", "delete_document_folder", "get_photo_folder_by_project_id", "create_photo_folder",
		"update_photo_folder", "delete_photo_folder", "get_chars_by_project_id", "create_char", "update_char", "delete_char",
		"get_states_by_project_id", "create_state", "update_state", "delete_state", "create_state_relation", "delete_state_relation",
		"get_projects_by_worker_id", "get_projects_by_client_id", "create_detail", "update_detail", "get_detail_by_id", "delete_detail",
	},
}

func Open(path string) (*gorm.DB, error) {
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Role{},
		&entity.UserRole{},
		&entity.Permission{},
		&entity.RolePermission{},
		&entity.Document{},
		&entity.Project{},
		&entity.ProjectRole{},
		&entity.DocumentFolder{},
		&entity.PhotoFolder{},
		&entity.Photo{},
		&entity.Char{},
		&entity.State{},
		&entity.StateUser{},
		&entity.RepairDetails{},
	)
	if err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

func SeedRoles(db *gorm.DB) error {
	all := map[string]bool{PermManageUsers: true, PermManageRoles: true}
	for _, perms := range rolePermissions {
		for _, p := range perms {
			all[p] = true
		}
	}

	permIDs := map[string]uint{}
	for name := range all {
		perm := entity.Permission{Name: name}
		if err := db.Where(entity.Permission{Name: name}).FirstOrCreate(&perm).Error; err != nil {
			return err
		}
		permIDs[name] = perm.ID
	}

	roles := map[string][]string{}
	for role, perms := range rolePermissions {
		roles[role] = perms
	}
	for name := range all {
		roles["admin"] = append(roles["admin"], name)
	}

	for roleName, perms := range roles {
		role := entity.Role{Name: roleName}
		if err := db.Where(entity.Role{Name: roleName}).FirstOrCreate(&role).Error; err != nil {
			return err
		}
		for _, p := range perms {
			link := entity.RolePermission{RoleID: role.ID, PermissionID: permIDs[p]}
			if err := db.Where(link).FirstOrCreate(&link).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func SeedUser(db *gorm.DB, user entity.User, password, roleName string) error {
	err := db.Where("email = ?", user.Email).First(&entity.User{}).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)

	var role entity.Role
	if err := db.Where("name = ?", roleName).First(&role).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return tx.Create(&entity.UserRole{UserID: user.ID, RoleID: role.ID}).Error
	})
}

func SeedDemoUsers(db *gorm.DB, password string) error {
	for _, role := range []string{"worker", "manager", "client", "prorab"} {
		user := entity.User{
			FirstName: role,
			LastName:  "Demo",
			Email:     role + "@demo.local",
		}
		if err := SeedUser(db, user, password, role); err != nil {
			return fmt.Errorf("seed %s: %w", role, err)
		}
	}
	return nil
}
