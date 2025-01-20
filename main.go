package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/HowkaCoder/remont/internal"
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/handler"
	"github.com/HowkaCoder/remont/internal/app/repository"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

var db *gorm.DB
var jwtSecret = []byte("your_secret_key")

func main() {
	db = internal.Init()

	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	docRepo := repository.NewDocumentRepository(db)
	docUsecase := usecase.NewDocumentUsecase(docRepo)
	docHandler := handler.NewDocumentHandler(docUsecase)

	projectRepo := repository.NewProjectRepository(db)
	projectUsecase := usecase.NewProjectUsecase(projectRepo)
	projectHandler := handler.NewProjectHandler(projectUsecase)

	photoRepo := repository.NewPhotoRepository(db)
	photoUsecase := usecase.NewPhotoUsecase(photoRepo)
	photoHandler := handler.NewPhotoHandler(photoUsecase)

	DocumentFolderRepo := repository.NewDocumentFolderRepository(db)
	documentFolderUsecase := usecase.NewDocumentFolderUsecase(DocumentFolderRepo)
	documentFolderHandler := handler.NewDocumentFolderHandler(documentFolderUsecase)

	CharRepo := repository.NewCharRepository(db)
	charUsecase := usecase.NewCharUsecase(CharRepo)
	charHandler := handler.NewCharHandler(charUsecase)

	photoFolder := repository.NewPhotoFolderRepository(db)
	photoFolderUsecase := usecase.NewPhotoFolderUsecase(photoFolder)
	photoFolderHandler := handler.NewPhotoFolderHandler(photoFolderUsecase)

	stateRepo := repository.NewStateRepository(db)
	stateUSecase := usecase.NewStateUsecase(stateRepo)
	stateHandler := handler.NewStateHandler(stateUSecase)

	app := fiber.New(fiber.Config{
		BodyLimit: 100 * 1024 * 1024,
	})
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE , PATCH")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Set("Access-Control-Allow-Credentials", "true")
		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Ошибка при получении текущей директории:", err)
		return

	}
	imagesDir := filepath.Join(currentDir, "uploads/photos/")

	documentDir := filepath.Join(currentDir, "uploads/documents/")
	// Статический обработчик для папки с изображениями
	app.Static("/uploads/photos", imagesDir)

	app.Static("/uploads/documents", documentDir)

	app.Get("/get-profile", ProtectedRoute)
	app.Delete("/users/:id", DeleteUser)
	// Роуты для регистрации и логина
	app.Post("/register", userHandler.CreateUser)
	app.Post("/login", login)
	app.Delete("/boom", clearDatabase)
	app.Post("/roles", createRole)
	app.Get("/roles", getAllRoles)
	app.Post("/permissions", createPermission)
	app.Get("/permissions", getAllPermissions)
	app.Post("/roles/assign", assignPermissionToRole)
	app.Post("/users/assign-role", assignRoleToUser)
	app.Get("/users", userHandler.GetAllUsers)
	app.Get("/api/projects/users/:id", getUsersByProjectID)
	app.Use(authMiddleware)

	app.Get("/api/docs", PermissionMiddleware("get_documents"), docHandler.GetAllDocuments)
	app.Post("/api/docs", PermissionMiddleware("create_document"), docHandler.CreateDocument)
	app.Patch("/api/docs/:id", PermissionMiddleware("update_document"), docHandler.UpdateDocument)
	app.Get("/api/docs/:id", PermissionMiddleware("get_document_by_id"), docHandler.GetDocumentByID)
	app.Get("/api/docs/folder/:folder_id", PermissionMiddleware("get_documents_by_folder_id"), docHandler.GetDocumentsByFolderID)
	app.Delete("/api/docs/:id", PermissionMiddleware("delete_document"), docHandler.DeleteDocument)

	app.Get("/api/projects", PermissionMiddleware("get_all_projects"), projectHandler.GetAllProjects)
	app.Post("/api/projects", PermissionMiddleware("create_project"), projectHandler.CreateProject)
	app.Patch("/api/projects/:id", PermissionMiddleware("update_project"), projectHandler.UpdateProject)
	app.Get("/api/projects/:id", PermissionMiddleware("get_project_by_id"), projectHandler.GetProjectByID)
	app.Delete("/api/projects/:id", PermissionMiddleware("delete_project"), projectHandler.DeleteProject)

	app.Get("/api/project/workers/:id", PermissionMiddleware("get_projects_by_worker_id"), projectHandler.GetAllProjectsAsAWorker)
	app.Get("/api/project/clients/:id", PermissionMiddleware("get_projects_by_client_id"), projectHandler.GetAllProjectsAsAClient)

	app.Get("/api/project-roles", PermissionMiddleware("get_all_project_roles"), projectHandler.GetAllProjectRole)
	app.Post("/api/project-roles", PermissionMiddleware("create_project_role"), projectHandler.CreateProjectRole)
	app.Patch("/api/project-roles/:id", PermissionMiddleware("update_project_role"), projectHandler.UpdateProjectRole)
	app.Get("/api/project-roles/:id", PermissionMiddleware("get_project_role_by_id"), projectHandler.GetProjectRoleByID)
	app.Delete("/api/project-roles/:id", PermissionMiddleware("delete_project_role"), projectHandler.DeleteProjectRole)

	app.Get("/api/photos", PermissionMiddleware("get_all_photos"), photoHandler.GetAllPhotos)
	app.Post("/api/photos", PermissionMiddleware("create_photo"), photoHandler.CreatePhoto)
	app.Get("/api/photos/folder/:folderID", PermissionMiddleware("get_photos_by_folder_id"), photoHandler.GetPhotosByFolderID)
	app.Get("/api/photos/:id", PermissionMiddleware("get_photo_by_id"), photoHandler.GetPhotoByID)
	app.Delete("/api/photos/:id", PermissionMiddleware("delete_photo"), photoHandler.DeletePhoto)
	app.Patch("/api/photos/:id", PermissionMiddleware("update_photo"), photoHandler.UpdatePhoto)

	app.Get("/api/document-folders", PermissionMiddleware("get_all_document_folders"), documentFolderHandler.GetAllDocumentFolders)
	app.Get("/api/document-folders/:id", PermissionMiddleware("get_document_folder_by_folder_id"), documentFolderHandler.GetDocumentFolderByID)
	app.Get("/api/document-folders/project/:id", PermissionMiddleware("get_document_folder_by_project_id"), documentFolderHandler.GetDocumentFoldersByProjectID)
	app.Post("/api/document-folders", PermissionMiddleware("create_document_folder"), documentFolderHandler.CreateDocumentFolder)
	app.Patch("/api/document-folders/:id", PermissionMiddleware("update_document_folder"), documentFolderHandler.UpdateDocumentFolder)
	app.Delete("/api/document-folders/:id", PermissionMiddleware("delete_document_folder"), documentFolderHandler.DeleteDocumentFolder)

	app.Get("/api/photo-folder/", PermissionMiddleware("get_all_photo_folders"), photoFolderHandler.GetAllPhotoFolders)
	app.Get("/api/photo-folder/:id", PermissionMiddleware("get_photo_folder_by_folder_id"), photoFolderHandler.GetPhotoFolderByID)
	app.Post("/api/photo-folder/", PermissionMiddleware("create_photo_folder"), photoFolderHandler.CreatePhotoFolder)
	app.Patch("/api/photo-folder/:id", PermissionMiddleware("update_photo_folder"), photoFolderHandler.UpdatePhotoFolder)
	app.Delete("/api/photo-folder/:id", PermissionMiddleware("delete_photo_folder"), photoFolderHandler.DeletePhotoFolder)
	app.Get("/api/photo-folder/project/:id", PermissionMiddleware("get_photo_folder_by_project_id"), photoFolderHandler.GetPhotoFoldersByProjectID)

	app.Get("/api/chars/project/:id", PermissionMiddleware("get_chars_by_project_id"), charHandler.GetAllCharsByProjectID)
	app.Post("/api/chars/project", PermissionMiddleware("create_char"), charHandler.CreateChar)
	app.Patch("/api/chars/:id", PermissionMiddleware("update_char"), charHandler.UpdateChar)
	app.Delete("/api/chars/:id", PermissionMiddleware("delete_char"), charHandler.DeleteChar)

	app.Get("/api/states/:id", PermissionMiddleware("get_states_by_project_id"), stateHandler.GetStatesByProjectID)
	app.Get("/api/states/worker/:id", PermissionMiddleware("get_states_by_worker_id"), stateHandler.GetStatesByWorkerID)
	app.Post("/api/states", PermissionMiddleware("create_state"), stateHandler.CreateState)
	app.Patch("/api/states/:id", PermissionMiddleware("update_state"), stateHandler.UpdateState)
	app.Delete("/api/states/:id", PermissionMiddleware("delete_state"), stateHandler.DeleteState)
	app.Post("/api/states/add-user", PermissionMiddleware("create_state_relation"), stateHandler.AssignWorkerToState)
	app.Post("/api/states/remove-user", PermissionMiddleware("delete_state_relation"), stateHandler.RemoveWorkerFromState)

	app.Post("/api/details", PermissionMiddleware("create_detail"), stateHandler.CreateRepairDetails)
	app.Get("/api/details/project/:id", PermissionMiddleware("get_detail_by_id"), stateHandler.GetRepairDetailsByProjectID)
	app.Patch("/api/details", PermissionMiddleware("update_detail"), stateHandler.UpdateRepairDetail)
	app.Delete("/api/details/:id", PermissionMiddleware("delete_detail"), stateHandler.DeleteRepairDetail)

	log.Fatal(app.Listen(":3000"))
}

func DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}

	var user entity.User

	db.First(&user, uint(id))

	db.Delete(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "successfully deleted"})
}

func PermissionMiddleware(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uint)
		var user entity.User
		db.Preload("Roles.Permissions").First(&user, userID)

		for _, role := range user.Roles {
			for _, permission := range role.Permissions {
				if permission.Name == requiredPermission {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}
}

/*
func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
	}

	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
	}
	claims, ok := token.Claims.(*entity.JwtClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	userID := claims.UserID
	if userID == 0 {
		return errors.New("userID is valid")
	}

	c.Locals("userID", userID)

	return c.Next()
}
*/

func getUsersByProjectID(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	var projectRoles []entity.ProjectRole
	db.Where("project_id = ?", uint(id)).Find(&projectRoles)

	var workerID []uint

	var manager uint
	for _, value := range projectRoles {
		if value.RoleID == 1 {
			workerID = append(workerID, value.UserID)
		}

		if value.RoleID == 2 {
			manager = value.UserID
		}
	}

	var users []entity.User

	for _, value := range workerID {
		var user entity.User

		db.First(&user, value)
		users = append(users, user)
	}

	var user entity.User
	db.First(&user, manager)

	return c.JSON(fiber.Map{
		"Workers": users,
		"Manager": user,
	})

}

func login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user entity.User
	db.Where("email = ?", req.Email).First(&user)
	if user.ID == 0 || user.Password != req.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	var userRole entity.UserRole
	db.Where("user_id = ?", user.ID).First(&userRole)

	var role entity.Role
	db.First(&role, userRole.RoleID)

	var projectRole []entity.ProjectRole
	db.Where("user_id = ?", user.ID).Find(&projectRole)

	var projects []uint
	for _, value := range projectRole {
		projects = append(projects, value.ProjectID)
	}
	claims := entity.JwtClaims{
		UserID:    user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      role.Name,
		Projects:  projects,
		City:      user.City,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Токен истекает через 24 часа
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": tokenString, "claims": claims})
}

func createRole(c *fiber.Ctx) error {
	var role entity.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	db.Create(&role)
	return c.Status(fiber.StatusCreated).JSON(role)
}

func createPermission(c *fiber.Ctx) error {
	var permission entity.Permission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	db.Create(&permission)
	return c.Status(fiber.StatusCreated).JSON(permission)
}

func assignPermissionToRole(c *fiber.Ctx) error {
	var input struct {
		RoleID       uint `json:"RoleID"`
		PermissionID uint `json:"PermissionID"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	rolePermission := entity.RolePermission{
		RoleID:       input.RoleID,
		PermissionID: input.PermissionID,
	}

	db.Create(&rolePermission)
	return c.Status(fiber.StatusCreated).JSON(rolePermission)
}

func assignRoleToUser(c *fiber.Ctx) error {
	var input struct {
		UserID uint `json:"UserID"`
		RoleID uint `json:"RoleID"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userRole := entity.UserRole{
		UserID: input.UserID,
		RoleID: input.RoleID,
	}

	db.Create(&userRole)
	return c.Status(fiber.StatusCreated).JSON(userRole)
}

func getAllRoles(c *fiber.Ctx) error {
	var roles []entity.Role
	if err := db.Preload("Permissions").Find(&roles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	return c.JSON(roles)
}

func getAllPermissions(c *fiber.Ctx) error {
	var permissions []entity.Permission
	if err := db.Find(&permissions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}
	return c.JSON(permissions)
}

func clearDatabase(c *fiber.Ctx) error {
	if _, err := os.Stat("remont.db"); err == nil {
		os.Remove("remont.db")
	}

	db = internal.Init()

	return c.JSON(fiber.Map{
		"message": "done",
	})
}
func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
	}

	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.ParseWithClaims(tokenString, &entity.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
	}

	claims, ok := token.Claims.(*entity.JwtClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	userID := claims.UserID
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	c.Locals("userID", userID)

	return c.Next()
}
func AuthenticateToken(tokenString string) (*entity.JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*entity.JwtClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
func ProtectedRoute(c *fiber.Ctx) error {
	tokenn := c.Get("Authorization")
	if tokenn == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed token"})
	}

	token := strings.TrimPrefix(tokenn, "Bearer ")
	claims, err := AuthenticateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	return c.JSON(fiber.Map{
		"user_id":    claims.UserID,
		"first_name": claims.FirstName,
		"last_name":  claims.LastName,
		"role":       claims.Role,
		"city":       claims.City,
	})
}

/*
func clearDatabase(c *fiber.Ctx) error {

	err := db.Migrator().DropTable(
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
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear database",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Database cleared successfully",
	})

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
		db.Create(&role)
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
	}

	for _, permission := range permissions {
		db.Create(&permission)
	}

	// Связывание ролей и разрешений
	rolePermissions := map[string][]string{
		"worker": {
			"get_all_projects", "create_project", "update_project", "get_project_by_id", "delete_project",
			"get_all_project_roles", "create_project_role", "update_project_role", "get_project_role_by_id",
			"delete_project_role", "get_all_photos", "create_photo", "get_photos_by_folder_id", "get_photo_by_id",
			"delete_photo", "update_photo", "get_states_by_worker_id",
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
			"get_projects_by_worker_id", "get_projects_by_client_id",
		},
		"client": {
			"get_documents", "get_document_by_id", "get_documents_by_folder_id", "get_all_projects", "get_project_by_id",
			"get_all_project_roles", "get_project_role_by_id", "get_all_photos", "get_photos_by_folder_id", "get_photo_by_id",
			"get_all_document_folders", "get_document_folder_by_folder_id", "get_document_folder_by_project_id",
			"get_all_photo_folders", "get_photo_folder_by_folder_id", "get_photo_folder_by_project_id",
			"get_chars_by_project_id", "get_states_by_project_id", "get_projects_by_worker_id", "get_projects_by_client_id",
		},
	}

	for roleName, perms := range rolePermissions {
		var role entity.Role
		db.Where("name = ?", roleName).First(&role)
		for _, permName := range perms {
			var perm entity.Permission
			db.Where("name = ?", permName).First(&perm)
			db.Create(&entity.RolePermission{
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
		db.Create(&user)
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
		db.Where("email = ?", email).First(&user)
		db.Where("name = ?", roleName).First(&role)
		db.Create(&entity.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		})
	}

	return c.JSON(fiber.Map{
		"message": "database successfully droped",
	})
} */
