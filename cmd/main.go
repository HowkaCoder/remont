package main

import (
	"log"

	"github.com/HowkaCoder/remont/internal"
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/handler"
	"github.com/HowkaCoder/remont/internal/app/repository"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
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

	app := fiber.New()

	// Роуты для регистрации и логина
	app.Post("/register", userHandler.CreateUser)
	app.Post("/login", login)

	app.Post("/roles", createRole)
	app.Get("/roles", getAllRoles)
	app.Post("/permissions", createPermission)
	app.Get("/permissions", getAllPermissions)
	app.Post("/roles/assign", assignPermissionToRole)
	app.Post("/users/assign-role", assignRoleToUser)
	app.Get("/users", userHandler.GetAllUsers)
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

	log.Fatal(app.Listen(":3000"))
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

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userID"].(float64))
	c.Locals("userID", userID)

	return c.Next()
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": tokenString})
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
	if err := db.Find(&roles).Error; err != nil {
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
