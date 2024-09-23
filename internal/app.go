package internal

import (
	"log"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	var err error
	//dsn := "host=dpg-cr7j2sjv2p9s73a556b0-a.oregon-postgres.render.com user=root password=OumUpk50PxWzWAu6Hni07HHvmdNj9SzE dbname=remont port=5432"
	DB, err = /* gorm.Open(postgres.Open(dsn), &gorm.Config{}) */ gorm.Open(sqlite.Open("database/database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Автоматическая миграция схемы базы данных
	err = DB.AutoMigrate(
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
		log.Fatal(err)
	}

	// Создание ролей
	roles := []entity.Role{
		{Name: "worker"},
		{Name: "manager"},
		{Name: "client"},
	}
	for _, role := range roles {
		DB.Create(&role)
	}

	// Создание разрешений
	permissions := []entity.Permission{
		{Name: "get_documents"},
		{Name: "create_document"},
		{Name: "update_document"},
		{Name: "get_document_by_id"},
		{Name: "get_documents_by_folder_id"},
		{Name: "delete_document"},
		{Name: "get_all_projects"},
		{Name: "create_project"},
		{Name: "update_project"},
		{Name: "get_project_by_id"},
		{Name: "delete_project"},
		{Name: "get_all_project_roles"},
		{Name: "create_project_role"},
		{Name: "update_project_role"},
		{Name: "get_project_role_by_id"},
		{Name: "delete_project_role"},
		{Name: "get_all_photos"},
		{Name: "create_photo"},
		{Name: "get_photos_by_folder_id"},
		{Name: "get_photo_by_id"},
		{Name: "delete_photo"},
		{Name: "update_photo"},
		{Name: "get_all_document_folders"},
		{Name: "get_document_folder_by_folder_id"},
		{Name: "get_document_folder_by_project_id"},
		{Name: "create_document_folder"},
		{Name: "update_document_folder"},
		{Name: "delete_document_folder"},
		{Name: "get_all_photo_folders"},
		{Name: "get_photo_folder_by_folder_id"},
		{Name: "create_photo_folder"},
		{Name: "update_photo_folder"},
		{Name: "delete_photo_folder"},
		{Name: "get_photo_folder_by_project_id"},
		{Name: "get_chars_by_project_id"},
		{Name: "create_char"},
		{Name: "update_char"},
		{Name: "delete_char"},
		{Name: "get_states_by_project_id"},
		{Name: "create_state"},
		{Name: "update_state"},
		{Name: "delete_state"},
		{Name: "create_state_relation"},
		{Name: "delete_state_relation"},
		{Name: "get_projects_by_worker_id"},
		{Name: "get_projects_by_client_id"},
		{Name: "get_states_by_worker_id"},
		{Name: "create_detail"},
		{Name: "update_detail"},
		{Name: "get_detail_by_id"},
		{Name: "delete_detail"},
	}

	for _, permission := range permissions {
		DB.Create(&permission)
	}

	// Связывание ролей и разрешений
	rolePermissions := map[string][]string{
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
	}

	for roleName, perms := range rolePermissions {
		var role entity.Role
		DB.Where("name = ?", roleName).First(&role)
		for _, permName := range perms {
			var perm entity.Permission
			DB.Where("name = ?", permName).First(&perm)
			DB.Create(&entity.RolePermission{
				RoleID:       role.ID,
				PermissionID: perm.ID,
			})
		}
	}

	// Создание пользователей
	users := []entity.User{
		{
			FirstName:   "worker",
			LastName:    "Khojaniyazov",
			MiddleName:  "Polatovich",
			Email:       "worker@gmail.com",
			Password:    "worker",
			PhoneNumber: "32e23e23e23e2",
		},
		{
			FirstName:   "manager",
			LastName:    "dewiojdwiejdoiwev",
			MiddleName:  "Polatovich",
			Email:       "manager@gmail.com",
			Password:    "manager",
			PhoneNumber: "32e23e23e23e2",
		},
		{
			FirstName:   "client",
			LastName:    "Khojaniyazov",
			MiddleName:  "Polatovich",
			Email:       "client@gmail.com",
			Password:    "client",
			PhoneNumber: "32e23e23e23e2",
		},
	}

	for _, user := range users {
		DB.Create(&user)
	}

	// Связывание пользователей и ролей
	userRoles := map[string]string{
		"worker@gmail.com":  "worker",
		"manager@gmail.com": "manager",
		"client@gmail.com":  "client",
	}

	for email, roleName := range userRoles {
		var user entity.User
		var role entity.Role
		DB.Where("email = ?", email).First(&user)
		DB.Where("name = ?", roleName).First(&role)
		DB.Create(&entity.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		})
	}

	return DB
}
