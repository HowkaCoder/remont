package internal

import (
	"log"

	"gorm.io/driver/sqlite"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	var err error
	//	dsn := "root:AvmOCFLHdwIkOcWYyXzGhuDvuTToYjsM@tcp(viaduct.proxy.rlwy.net:38909)/railway?charset=utf8mb4&parseTime=True&loc=Local"
	//DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	DB, err := gorm.Open(sqlite.Open("database/database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

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
		&entity.ProjectRole{},
		&entity.Char{},
		&entity.State{},
		&entity.StateUser{},
		//		&entity.ProjectManager{},
	//	&entity.ProjectWorker{},
	)
	if err != nil {
		log.Fatal(err)
	}

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
	}

	for _, permission := range permissions {
		DB.Create(&permission)

		rolePermission := entity.RolePermission{
			RoleID:       1,
			PermissionID: permission.ID,
		}

		DB.Create(&rolePermission)
	}

	return DB

}
