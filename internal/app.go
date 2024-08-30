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
	// postgresql://root:OumUpk50PxWzWAu6Hni07HHvmdNj9SzE@dpg-cr7j2sjv2p9s73a556b0-a/remont
	//	dsn := "root:AvmOCFLHdwIkOcWYyXzGhuDvuTToYjsM@tcp(viaduct.proxy.rlwy.net:38909)/railway?charset=utf8mb4&parseTime=True&loc=Local"
	//postgresql://root:OumUpk50PxWzWAu6Hni07HHvmdNj9SzE@dpg-cr7j2sjv2p9s73a556b0-a.oregon-postgres.render.com/remont/
	//DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//	dsn := "host=dpg-cr7j2sjv2p9s73a556b0-a.oregon-postgres.render.com user=root password=OumUpk50PxWzWAu6Hni07HHvmdNj9SzE dbname=remont port=5432 "
	//	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB, err := gorm.Open(sqlite.Open("./database/database.db"), &gorm.Config{})
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

	/*	role1 := entity.Role{
			Name: "manager",
		}

		DB.Create(role1)

		role2 := entity.Role{
			Name: "director",
		}

		DB.Create(role2)

		role3 := entity.Role{
			Name: "client",
		}
		DB.Create(role3)

		role4 := entity.Role{
			Name: "worker",
		}

		DB.Create(role4)*/

	roles := []entity.Role{
		{Name: "worker"},
		{Name: "manager"},
		{Name: "client"},
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
		{Name: "get_projects_by_worker_id"},
		{Name: "get_projects_by_client_id"},
		//	{Name: "get_states_by_worker_id"},
	}

	user1 := entity.User{
		FirstName:   "worker",
		LastName:    "Khojaniyazov",
		MiddleName:  "Polatovich",
		Email:       "worker@gmail.com",
		Password:    "worker",
		PhoneNumber: "32e23e23e23e2",
	}
	DB.Create(&user1)

	user2 := entity.User{
		FirstName:   "manager",
		LastName:    "dewiojdwiejdoiwev",
		MiddleName:  "Polatovich",
		Email:       "manager@gmail.com",
		Password:    "manager",
		PhoneNumber: "32e23e23e23e2",
	}
	DB.Create(&user2)

	user3 := entity.User{
		FirstName:   "client",
		LastName:    "Khojaniyazov",
		MiddleName:  "Polatovich",
		Email:       "client@gmail.com",
		Password:    "client",
		PhoneNumber: "32e23e23e23e2",
	}
	DB.Create(&user3)

	for _, permission := range permissions {
		DB.Create(&permission)

		rolePermission := entity.RolePermission{
			RoleID:       2,
			PermissionID: permission.ID,
		}

		DB.Create(&rolePermission)
		/*
			rolePermission1 := entity.RolePermission{
				RoleID:       2,
				PermissionID: permission.ID,
			}

			DB.Create(&rolePermission1)

			rolePermission3 := entity.RolePermission{
				RoleID:       3,
				PermissionID: permission.ID,
			}

			DB.Create(&rolePermission3)

			rolePermission2 := entity.RolePermission{
				RoleID:       4,
				PermissionID: permission.ID,
			}

			DB.Create(&rolePermission2)
		*/

	}

	for _, role := range roles {
		DB.Create(&role)

		userRole := entity.UserRole{
			UserID: role.ID,
			RoleID: role.ID,
		}

		DB.Create(&userRole)
	}

	WorkerRolePermissions := []entity.RolePermission{
		{
			RoleID:       1,
			PermissionID: 17,
		},
		{
			RoleID:       1,
			PermissionID: 18,
		},
		{
			RoleID:       1,
			PermissionID: 19,
		},

		{
			RoleID:       1,
			PermissionID: 20,
		},
		{
			RoleID:       1,
			PermissionID: 21,
		},
		{
			RoleID:       1,
			PermissionID: 22,
		},
		{
			RoleID:       1,
			PermissionID: 36,
		},
		{
			RoleID:       1,
			PermissionID: 37,
		},
		{
			RoleID:       1,
			PermissionID: 38,
		},

		{
			RoleID:       1,
			PermissionID: 39,
		},

		{
			RoleID:       1,
			PermissionID: 40,
		},

		{
			RoleID:       1,
			PermissionID: 41,
		},

		{
			RoleID:       1,
			PermissionID: 42,
		},
		{
			RoleID:       1,
			PermissionID: 45,
		},
	}

	for _, workerRolePermission := range WorkerRolePermissions {
		DB.Create(&workerRolePermission)
	}

	ClientRolePermissions := []entity.RolePermission{
		{
			RoleID:       3,
			PermissionID: 1,
		},
		{
			RoleID:       3,
			PermissionID: 4,
		},
		{
			RoleID:       3,
			PermissionID: 5,
		},
		{
			RoleID:       3,
			PermissionID: 17,
		},
		{
			RoleID:       3,
			PermissionID: 46,
		},
		{
			RoleID:       3,
			PermissionID: 47,
		},
		{
			RoleID:       3,
			PermissionID: 18,
		},
		{
			RoleID:       3,
			PermissionID: 19,
		},
		{
			RoleID:       3,
			PermissionID: 20,
		},
		{
			RoleID:       3,
			PermissionID: 21,
		},
		{
			RoleID:       3,
			PermissionID: 22,
		},
		{
			RoleID:       3,
			PermissionID: 39,
		},
		{
			RoleID:       3,
			PermissionID: 40,
		},
		{
			RoleID:       3,
			PermissionID: 41,
		},
	}
	for _, clientRolePermission := range ClientRolePermissions {
		DB.Create(clientRolePermission)
	}

	DB.Create(&entity.Permission{Name: "get_states_by_worker_id"})
	return DB

}
